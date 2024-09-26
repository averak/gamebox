package echo_usecase

import (
	"context"
	"fmt"

	"github.com/averak/gamebox/app/core/game_context"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository/transaction"
)

func (u Usecase) Echo(ctx context.Context, gctx game_context.GameContext, message string) (model.Echo, error) {
	echo := model.NewEcho(message, gctx.Now())
	err := u.conn.BeginRwTransaction(ctx, func(ctx context.Context, tx transaction.Transaction) error {
		err := u.echoRepo.Save(ctx, tx, echo)
		if err != nil {
			return fmt.Errorf("echoRepo.Save failed, %w", err)
		}
		return nil
	})
	if err != nil {
		return model.Echo{}, err
	}
	return echo, nil
}
