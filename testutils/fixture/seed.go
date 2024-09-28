package fixture

import (
	"context"
	"testing"

	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/testutils"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Seed interface {
	Insert(context.Context, boil.ContextExecutor, boil.Columns) error
}

func SetupSeeds(t *testing.T, ctx context.Context, seeds ...Seed) {
	t.Helper()

	conn := testutils.MustDBConn(t)
	err := conn.BeginRwTransaction(ctx, func(ctx context.Context, tx transaction.Transaction) error {
		// 外部キー制約を守るために、非効率だが直列実行する必要がある。
		// ユーティリティ側で依存を解決できると嬉しいが、難しいので直列実行で妥協する。
		for _, seed := range seeds {
			if err := seed.Insert(ctx, tx, boil.Infer()); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
