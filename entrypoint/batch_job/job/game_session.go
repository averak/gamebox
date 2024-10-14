package job

import (
	"context"
	"time"

	"github.com/averak/gamebox/app/adapter/dao"
	"github.com/averak/gamebox/app/core/game_context"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/volatiletech/null/v8"
)

type purgeOldGameSessionsJob struct{}

func NewPurgeOldGameSessions() BatchJob {
	return purgeOldGameSessionsJob{}
}

func (p purgeOldGameSessionsJob) Desc() string {
	return "プレイ後 90d 以上経過したゲームセッションを削除します。"
}

func (p purgeOldGameSessionsJob) Run(ctx context.Context, gctx game_context.GameContext, conn transaction.Connection) error {
	ttl := 90 * 24 * time.Hour
	// 現在時刻が 10/5 だとすると、9/1 00:00:00 よりも前のレコードを削除する
	return conn.BeginRwTransaction(ctx, func(ctx context.Context, tx transaction.Transaction) error {
		_, err := dao.UserGameSessions(dao.UserGameSessionWhere.FinishedAt.LTE(null.NewTime(gctx.Now().Add(-ttl), true))).DeleteAll(ctx, tx)
		return err
	})
}
