package game_session_repoimpl

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
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Repository struct{}

func NewRepository() repository.GameSessionRepository {
	return &Repository{}
}

func (r Repository) Get(ctx context.Context, tx transaction.Transaction, gameSessionID uuid.UUID, userID uuid.UUID) (model.GameSession, error) {
	ctx, span := trace.StartSpan(ctx, "game_session_repoimpl.Get")
	defer span.End()

	dto, err := dao.UserGameSessions(
		dao.UserGameSessionWhere.UserID.EQ(userID.String()),
		dao.UserGameSessionWhere.ID.EQ(gameSessionID.String()),
	).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.GameSession{}, repository.ErrGameSessionNotFound
		}
		return model.GameSession{}, err
	}
	wager, err := model.NewCoins(dto.Wager)
	if err != nil {
		return model.GameSession{}, err
	}
	payout, err := model.NewCoins(dto.Payout)
	if err != nil {
		return model.GameSession{}, err
	}
	return model.NewGameSession(
		uuid.MustParse(dto.ID),
		uuid.MustParse(dto.UserID),
		model.GameID(dto.GameID),
		model.GameStatus(dto.Status),
		model.GameResult(dto.Result),
		wager,
		payout,
		dto.StartedAt,
		dto.FinishedAt.Time,
	)
}

func (r Repository) Save(ctx context.Context, tx transaction.Transaction, sessions ...model.GameSession) error {
	ctx, span := trace.StartSpan(ctx, "game_session_repoimpl.Save")
	defer span.End()

	dtos := vector.Map(sessions, func(session model.GameSession) *dao.UserGameSession {
		return &dao.UserGameSession{
			ID:         session.ID.String(),
			UserID:     session.UserID.String(),
			GameID:     int(session.GameID),
			Status:     int(session.Status),
			Result:     int(session.Result),
			Wager:      int(session.Wager),
			Payout:     int(session.Payout),
			StartedAt:  session.StartedAt,
			FinishedAt: null.NewTime(session.FinishedAt, !session.FinishedAt.IsZero()),
		}
	})
	_, err := dao.UserGameSessionSlice(dtos).UpsertAll(ctx, tx, true, dao.UserGameSessionPrimaryKeyColumns, boil.Whitelist(dao.UserGameSessionAllColumns...), boil.Whitelist(dao.UserGameSessionAllColumns...))
	if err != nil {
		return err
	}
	return nil
}
