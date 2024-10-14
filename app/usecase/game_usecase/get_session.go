package game_usecase

import (
	"context"

	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/google/uuid"
)

func (u Usecase) GetSession(ctx context.Context, user model.User, gameSessionID uuid.UUID) (model.GameSession, error) {
	var res model.GameSession
	err := u.conn.BeginRoTransaction(ctx, func(ctx context.Context, tx transaction.Transaction) error {
		var err error
		res, err = u.gameSessionRepo.Get(ctx, tx, gameSessionID, user.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.GameSession{}, err
	}
	return res, nil
}
