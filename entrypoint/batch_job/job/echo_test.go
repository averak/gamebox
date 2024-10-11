package job

import (
	"context"
	"testing"
	"time"

	"github.com/averak/gamebox/app/adapter/dao"
	"github.com/averak/gamebox/app/core/game_context"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/pkg/vector"
	"github.com/averak/gamebox/testutils"
	"github.com/averak/gamebox/testutils/bdd"
	"github.com/averak/gamebox/testutils/faker"
	"github.com/averak/gamebox/testutils/fixture"
	"github.com/stretchr/testify/require"
)

func TestPurgeOldEchos_Run(t *testing.T) {
	conn := testutils.MustDBConn(t)
	now := time.Now()

	type given struct {
		seeds []fixture.Seed
	}
	type when struct {
		gctx game_context.GameContext
	}
	type then = func(*testing.T, []*dao.Echo, error)
	tests := []bdd.Testcase[given, when, then]{
		{
			Given: given{
				seeds: []fixture.Seed{
					&dao.Echo{
						ID:        faker.UUIDv5("e1").String(),
						CreatedAt: now.Add(-23 * time.Hour),
					},
					&dao.Echo{
						ID:        faker.UUIDv5("e2").String(),
						CreatedAt: now.Add(-24 * time.Hour),
					},
					&dao.Echo{
						ID:        faker.UUIDv5("e3").String(),
						CreatedAt: now.Add(-25 * time.Hour),
					},
					&dao.Echo{
						ID:        faker.UUIDv5("e4").String(),
						CreatedAt: now.Add(-26 * time.Hour),
					},
				},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "24時間前に作成されたレコードが削除される",
					When: when{
						gctx: faker.NewGameContextBuilder().Now(now).Build(),
					},
					Then: func(t *testing.T, dtos []*dao.Echo, err error) {
						require.NoError(t, err)
						gotIDs := vector.Map(dtos, func(dto *dao.Echo) string { return dto.ID })
						require.Equal(t, []string{faker.UUIDv5("e1").String(), faker.UUIDv5("e2").String()}, gotIDs)
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.Run(t, func(t *testing.T, given given, when when, then then) {
			fixture.SetupSeeds(t, context.Background(), given.seeds...)
			defer testutils.Teardown(t)

			j := purgeOldEchosJob{}
			err := j.Run(context.Background(), when.gctx, conn)

			var dtos []*dao.Echo
			txErr := conn.BeginRoTransaction(context.Background(), func(ctx context.Context, tx transaction.Transaction) error {
				var err error
				dtos, err = dao.Echos().All(ctx, tx)
				return err
			})
			if txErr != nil {
				t.Fatal(txErr)
			}
			then(t, dtos, err)
		})
	}
}
