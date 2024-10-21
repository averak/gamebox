package game_usecase

import (
	"context"

	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository/transaction"
)

func (u Usecase) ListPlayingSession(ctx context.Context, user model.User) ([]model.GameSession, error) {
	var res []model.GameSession
	err := u.conn.BeginRoTransaction(ctx, func(ctx context.Context, tx transaction.Transaction) error {
		var err error
		res, err = u.gameSessionRepo.GetByStatus(ctx, tx, user.ID, model.GameStatusPlaying)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
