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
	"github.com/averak/gamebox/testutils/fixture/fixture_builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_purgeOldGameSessionsJob_Run(t *testing.T) {
	conn := testutils.MustDBConn(t)
	now := time.Now()

	type given struct {
		fixtures []fixture_builder.Fixture
	}
	type when struct {
		gctx game_context.GameContext
	}
	type then = func(*testing.T, []*dao.UserGameSession, error)
	tests := []bdd.Testcase[given, when, then]{
		{
			Given: given{
				fixtures: []fixture_builder.Fixture{
					fixture_builder.New().
						User(fixture_builder.NewUserBuilder(faker.UUIDv5("u1")).Build()).
						GameSession(
							fixture_builder.NewGameSessionBuilder(t, faker.UUIDv5("u1"), faker.UUIDv5("gs1")).
								FinishedAt(now.Add(-90*24*time.Hour)).
								Build(),
							fixture_builder.NewGameSessionBuilder(t, faker.UUIDv5("u1"), faker.UUIDv5("gs2")).
								FinishedAt(now.Add(-90*24*time.Hour).Add(time.Millisecond)).
								Build(),
							fixture_builder.NewGameSessionBuilder(t, faker.UUIDv5("u1"), faker.UUIDv5("gs3")).
								FinishedAt(time.Time{}).
								Build(),
						).
						Build(),
				},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "プレイ後 90d 以上経過したゲームセッションを削除する",
					When: when{
						gctx: faker.NewGameContextBuilder().Now(now).Build(),
					},
					Then: func(t *testing.T, dtos []*dao.UserGameSession, err error) {
						require.NoError(t, err)

						gotIDs := vector.Map(dtos, func(dto *dao.UserGameSession) string { return dto.ID })
						assert.ElementsMatch(t, []string{faker.UUIDv5("gs2").String(), faker.UUIDv5("gs3").String()}, gotIDs)
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.Run(t, func(t *testing.T, given given, when when, then then) {
			fixture.Setup(t, given.fixtures...)
			defer testutils.Teardown(t)

			j := NewPurgeOldGameSessions()
			err := j.Run(context.Background(), when.gctx, conn)

			var dtos []*dao.UserGameSession
			txErr := conn.BeginRoTransaction(context.Background(), func(ctx context.Context, tx transaction.Transaction) error {
				var err error
				dtos, err = dao.UserGameSessions().All(ctx, tx)
				return err
			})
			if txErr != nil {
				t.Fatal(txErr)
			}
			then(t, dtos, err)
		})
	}
}
