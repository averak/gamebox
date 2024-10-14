package janken_session_repoimpl

import (
	"context"
	"database/sql"
	"errors"

	"github.com/averak/gamebox/app/adapter/dao"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/app/infrastructure/trace"
	"github.com/averak/gamebox/pkg/vector"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Repository struct{}

func NewRepository() repository.JankenSessionRepository {
	return &Repository{}
}

func (r Repository) Get(ctx context.Context, tx transaction.Transaction, gameSessionID uuid.UUID) (model.JankenSession, error) {
	ctx, span := trace.StartSpan(ctx, "janken_session_repoimpl.Get")
	defer span.End()

	dto, err := dao.UserJankenSessions(
		qm.Load(dao.UserJankenSessionRels.GameSessionUserJankenSessionHistories),
		dao.UserJankenSessionWhere.GameSessionID.EQ(gameSessionID.String()),
	).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.JankenSession{}, repository.ErrJankenSessionNotFound
		}
		return model.JankenSession{}, err
	}

	histories := vector.Map(dto.R.GameSessionUserJankenSessionHistories, func(dto *dao.UserJankenSessionHistory) model.JankenHistory {
		return model.NewJankenHistory(dto.Turn, model.JankenHand(dto.MyHand), model.JankenHand(dto.OpponentHand))
	})
	return model.NewJankenSession(gameSessionID, dto.Seed, histories), nil
}

func (r Repository) Save(ctx context.Context, tx transaction.Transaction, session model.JankenSession) error {
	ctx, span := trace.StartSpan(ctx, "janken_session_repoimpl.Save")
	defer span.End()
	boil.DebugMode = true

	{
		dto := &dao.UserJankenSession{
			GameSessionID: session.GameSessionID.String(),
			Seed:          session.Seed,
		}
		err := dto.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer())
		if err != nil {
			return err
		}
	}
	{
		dtos := vector.Map(session.Histories, func(history model.JankenHistory) *dao.UserJankenSessionHistory {
			return &dao.UserJankenSessionHistory{
				GameSessionID: session.GameSessionID.String(),
				Turn:          history.Turn,
				MyHand:        int(history.MyHand),
				OpponentHand:  int(history.OpponentHand),
			}
		})
		if len(dtos) == 0 {
			// 履歴が存在しない場合、後続の処理が無駄になるのでスキップする。
			return nil
		}

		currentDtos, err := dao.UserJankenSessionHistories(dao.UserJankenSessionHistoryWhere.GameSessionID.EQ(session.GameSessionID.String())).All(ctx, tx)
		if err != nil {
			return err
		}

		upserted, deleted := CheckHistoryDiff(dtos, currentDtos)
		_, err = dao.UserJankenSessionHistorySlice(upserted).UpsertAll(ctx, tx, true, dao.UserJankenSessionHistoryPrimaryKeyColumns, boil.Infer(), boil.Infer())
		if err != nil {
			return err
		}

		_, err = dao.UserJankenSessionHistorySlice(deleted).DeleteAll(ctx, tx)
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckHistoryDiff は、新しい履歴と現在の履歴を比較して、作成/更新/削除された履歴を仕分けます。
func CheckHistoryDiff(newDtos, currentDtos []*dao.UserJankenSessionHistory) (upserted []*dao.UserJankenSessionHistory, deleted []*dao.UserJankenSessionHistory) {
	currentMap := make(map[int]*dao.UserJankenSessionHistory)
	for _, current := range currentDtos {
		currentMap[current.Turn] = current
	}

	upserted = make([]*dao.UserJankenSessionHistory, 0)
	for _, dto := range newDtos {
		current, ok := currentMap[dto.Turn]
		if !ok {
			// 作成された
			upserted = append(upserted, dto)
		} else if !cmp.Equal(dto, current, cmpopts.IgnoreFields(dao.UserJankenSessionHistory{}, "CreatedAt", "UpdatedAt")) {
			// 更新された
			upserted = append(upserted, dto)
		}
		delete(currentMap, dto.Turn)
	}

	deleted = make([]*dao.UserJankenSessionHistory, 0)
	for _, leftover := range currentMap {
		deleted = append(deleted, leftover)
	}
	return upserted, deleted
}
