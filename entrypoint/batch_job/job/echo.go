package job

import (
	"context"
	"time"

	"github.com/averak/gamebox/app/adapter/dao"
	"github.com/averak/gamebox/app/core/game_context"
	"github.com/averak/gamebox/app/domain/repository/transaction"
)

type purgeOldEchosJob struct{}

func NewPurgeOldEchos() BatchJob {
	return purgeOldEchosJob{}
}

func (j purgeOldEchosJob) Desc() string {
	return "echos テーブルの古いレコードを削除します。"
}

func (j purgeOldEchosJob) Run(ctx context.Context, gctx game_context.GameContext, conn transaction.Connection) error {
	ttl := 24 * time.Hour
	return conn.BeginRwTransaction(ctx, func(ctx context.Context, tx transaction.Transaction) error {
		_, err := dao.Echos(dao.EchoWhere.CreatedAt.LT(gctx.Now().Add(-ttl))).DeleteAll(ctx, tx)
		return err
	})
}
