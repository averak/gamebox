package job

import (
	"context"

	"github.com/averak/gamebox/app/core/game_context"
	"github.com/averak/gamebox/app/domain/repository/transaction"
)

// BatchJob は、以下の要件を満たす必要があります。
// 1. リトライ可能であること (冪等性を保証する)。
// 2. 時間のかかる処理をレコードごとに実行しないこと (これが難しい場合は、非同期実行することを検討してください)。
type BatchJob interface {
	Desc() string
	Run(ctx context.Context, gctx game_context.GameContext, conn transaction.Connection) error
}
