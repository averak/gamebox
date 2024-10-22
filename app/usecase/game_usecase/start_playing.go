package game_usecase

import (
	"context"

	"github.com/averak/gamebox/app/core/game_context"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository/transaction"
)

func (u Usecase) StartPlaying(ctx context.Context, gctx game_context.GameContext, user model.User, gameID model.GameID, wager model.Coins) (model.GameSession, error) {
	var sess model.GameSession
	err := u.conn.BeginRwTransaction(ctx, func(ctx context.Context, tx transaction.Transaction) error {
		playingSession, err := u.gameSessionRepo.GetByStatus(ctx, tx, user.ID, model.GameStatusPlaying)
		if err != nil {
			return err
		}
		s, err := model.NewGameSessionService(ctx, user.ID, playingSession)
		if err != nil {
			return err
		}
		sess, err = s.StartPlaying(gctx.IdempotencyKey(), gameID, wager, gctx.Now())
		if err != nil {
			return err
		}
		err = u.gameSessionRepo.Save(ctx, tx, sess)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.GameSession{}, err
	}
	return sess, nil
}
