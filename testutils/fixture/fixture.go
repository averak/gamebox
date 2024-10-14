package fixture

import (
	"context"
	"testing"

	"github.com/averak/gamebox/app/adapter/repoimpl/game_session_repoimpl"
	"github.com/averak/gamebox/app/adapter/repoimpl/user_repoimpl"
	"github.com/averak/gamebox/app/domain/repository"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/testutils"
	"github.com/averak/gamebox/testutils/fixture/fixture_builder"
)

var repos = struct {
	UserRepo        repository.UserRepository
	GameSessionRepo repository.GameSessionRepository
}{
	UserRepo:        user_repoimpl.NewRepository(),
	GameSessionRepo: game_session_repoimpl.NewRepository(),
}

func Setup(t *testing.T, fixtures ...fixture_builder.Fixture) {
	t.Helper()

	if len(fixtures) == 0 {
		return
	}

	conn := testutils.MustDBConn(t)
	err := conn.BeginRwTransaction(context.Background(), func(ctx context.Context, tx transaction.Transaction) error {
		for _, f := range fixtures {
			for _, v := range f.Users {
				err := repos.UserRepo.Save(ctx, tx, v)
				if err != nil {
					return err
				}
			}

			err := repos.GameSessionRepo.Save(ctx, tx, f.GameSessions...)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
