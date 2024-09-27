package game_session_repoimpl

import (
	"context"
	"testing"
	"time"

	"github.com/averak/gamebox/app/adapter/dao"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/testutils"
	"github.com/averak/gamebox/testutils/bdd"
	"github.com/averak/gamebox/testutils/faker"
	"github.com/averak/gamebox/testutils/fixture"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
)

func TestRepository_Get(t *testing.T) {
	conn := testutils.MustDBConn(t)
	now := time.Now().Truncate(time.Millisecond)

	type given struct {
		seeds []fixture.Seed
	}
	type when struct {
		id uuid.UUID
	}
	type then = func(t *testing.T, got model.GameSession, err error)
	tests := []bdd.Testcase[given, when, then]{
		{
			Name: "レコードが存在する状態で",
			Given: given{
				seeds: []fixture.Seed{
					&dao.User{
						ID: faker.UUIDv5("u1").String(),
					},
					&dao.UserGameSession{
						ID:         faker.UUIDv5("gs1").String(),
						UserID:     faker.UUIDv5("u1").String(),
						GameID:     int(model.GameIDBlackjack),
						Status:     int(model.GameStatusPlaying),
						Result:     int(model.GameResultUnknown),
						Wager:      100,
						Payout:     0,
						StartedAt:  now,
						FinishedAt: null.Time{},
					},
				},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "ID で検索できる",
					When: when{
						id: faker.UUIDv5("gs1"),
					},
					Then: func(t *testing.T, got model.GameSession, err error) {
						require.NoError(t, err)

						want := model.GameSession{
							ID:         faker.UUIDv5("gs1"),
							UserID:     faker.UUIDv5("u1"),
							GameID:     model.GameIDBlackjack,
							Status:     model.GameStatusPlaying,
							Result:     model.GameResultUnknown,
							Wager:      100,
							Payout:     0,
							StartedAt:  now,
							FinishedAt: time.Time{},
						}
						assert.EqualExportedValues(t, want, got)
					},
				},
				{
					Name: "ID が存在しない => エラー",
					When: when{
						id: faker.UUIDv5("not found"),
					},
					Then: func(t *testing.T, got model.GameSession, err error) {
						assert.ErrorIs(t, err, repository.ErrGameSessionNotFound)
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.Run(t, func(t *testing.T, given given, when when, then then) {
			defer testutils.Teardown(t)
			fixture.SetupSeeds(t, context.Background(), given.seeds...)

			var got model.GameSession
			err := conn.BeginRoTransaction(context.Background(), func(ctx context.Context, tx transaction.Transaction) error {
				r := NewRepository()
				var err error
				got, err = r.Get(ctx, tx, when.id)
				return err
			})
			then(t, got, err)
		})
	}
}

func TestRepository_Save(t *testing.T) {
	conn := testutils.MustDBConn(t)
	now := time.Now().Truncate(time.Millisecond)

	type given struct {
		seeds []fixture.Seed
	}
	type when struct {
		sessions []model.GameSession
	}
	type then = func(t *testing.T, dtos []*dao.UserGameSession, err error)
	tests := []bdd.Testcase[given, when, then]{
		{
			Name: "レコードが存在する状態で",
			Given: given{
				seeds: []fixture.Seed{
					&dao.User{
						ID: faker.UUIDv5("u1").String(),
					},
					&dao.UserGameSession{
						ID:     faker.UUIDv5("gs1").String(),
						UserID: faker.UUIDv5("u1").String(),
						GameID: int(model.GameIDBlackjack),
					},
				},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "UPSERT できる",
					When: when{
						sessions: []model.GameSession{
							{
								ID:         faker.UUIDv5("gs1"),
								UserID:     faker.UUIDv5("u1"),
								GameID:     model.GameIDBlackjack,
								Status:     model.GameStatusPlaying,
								Result:     model.GameResultUnknown,
								Wager:      100,
								Payout:     0,
								StartedAt:  now,
								FinishedAt: time.Time{},
							},
							{
								ID:         faker.UUIDv5("gs2"),
								UserID:     faker.UUIDv5("u1"),
								GameID:     model.GameIDSolitaire,
								Status:     model.GameStatusFinished,
								Result:     model.GameResultWin,
								Wager:      100,
								Payout:     200,
								StartedAt:  now,
								FinishedAt: now,
							},
						},
					},
					Then: func(t *testing.T, dtos []*dao.UserGameSession, err error) {
						require.NoError(t, err)

						want := []*dao.UserGameSession{
							{
								ID:         faker.UUIDv5("gs1").String(),
								UserID:     faker.UUIDv5("u1").String(),
								GameID:     int(model.GameIDBlackjack),
								Status:     int(model.GameStatusPlaying),
								Result:     int(model.GameResultUnknown),
								Wager:      100,
								Payout:     0,
								StartedAt:  now,
								FinishedAt: null.Time{},
							},
							{
								ID:         faker.UUIDv5("gs2").String(),
								UserID:     faker.UUIDv5("u1").String(),
								GameID:     int(model.GameIDSolitaire),
								Status:     int(model.GameStatusFinished),
								Result:     int(model.GameResultWin),
								Wager:      100,
								Payout:     200,
								StartedAt:  now,
								FinishedAt: null.NewTime(now, true),
							},
						}
						if diff := cmp.Diff(want, dtos, cmpopts.IgnoreFields(dao.UserGameSession{}, "CreatedAt", "UpdatedAt")); diff != "" {
							t.Errorf("(-want, +got)\n%s", diff)
						}
					},
				},
				{
					Name: "空リスト => 何もしない",
					When: when{
						sessions: []model.GameSession{},
					},
					Then: func(t *testing.T, dtos []*dao.UserGameSession, err error) {
						require.NoError(t, err)

						want := []*dao.UserGameSession{
							{
								ID:     faker.UUIDv5("gs1").String(),
								UserID: faker.UUIDv5("u1").String(),
								GameID: int(model.GameIDBlackjack),
							},
						}
						if diff := cmp.Diff(want, dtos, cmpopts.IgnoreFields(dao.UserGameSession{}, "CreatedAt", "UpdatedAt")); diff != "" {
							t.Errorf("(-want, +got)\n%s", diff)
						}
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.Run(t, func(t *testing.T, given given, when when, then then) {
			defer testutils.Teardown(t)
			fixture.SetupSeeds(t, context.Background(), given.seeds...)

			var dtos []*dao.UserGameSession
			err := conn.BeginRwTransaction(context.Background(), func(ctx context.Context, tx transaction.Transaction) error {
				r := NewRepository()
				err := r.Save(ctx, tx, when.sessions...)
				if err != nil {
					return err
				}

				dtos, err = dao.UserGameSessions().All(ctx, tx)
				return err
			})
			then(t, dtos, err)
		})
	}
}
