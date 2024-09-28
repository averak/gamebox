package echo_repoimpl

import (
	"context"
	"testing"
	"time"

	"github.com/averak/gamebox/app/adapter/dao"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/testutils"
	"github.com/averak/gamebox/testutils/bdd"
	"github.com/averak/gamebox/testutils/faker"
	"github.com/averak/gamebox/testutils/fixture"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestRepository_Save(t *testing.T) {
	conn := testutils.MustDBConn(t)
	now := time.Now().Truncate(time.Millisecond)

	type given struct {
		seeds []fixture.Seed
	}
	type when struct {
		echos []model.Echo
	}
	type then = func(t *testing.T, dtos []*dao.Echo, err error)
	tests := []bdd.Testcase[given, when, then]{
		{
			Name: "レコードが存在する状態で",
			Given: given{
				seeds: []fixture.Seed{
					&dao.Echo{ID: faker.UUIDv5("e1").String(), Message: "m1", Timestamp: now},
				},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "PK が重複する場合は更新し、重複しない場合は作成される",
					When: when{
						echos: []model.Echo{
							{
								ID:        faker.UUIDv5("e1"),
								Message:   "updated",
								Timestamp: now.Add(1 * time.Hour),
							},
							{
								ID:        faker.UUIDv5("e2"),
								Message:   "created",
								Timestamp: now.Add(2 * time.Hour),
							},
						},
					},
					Then: func(t *testing.T, dtos []*dao.Echo, err error) {
						require.NoError(t, err)

						want := []*dao.Echo{
							{
								ID:        dtos[0].ID,
								Message:   "updated",
								Timestamp: now.Add(1 * time.Hour),
							},
							{
								ID:        dtos[1].ID,
								Message:   "created",
								Timestamp: now.Add(2 * time.Hour),
							},
						}
						if diff := cmp.Diff(want, dtos, cmpopts.IgnoreFields(dao.Echo{}, "CreatedAt", "UpdatedAt")); diff != "" {
							t.Errorf("(-want, +got)\n%s", diff)
						}
					},
				},
				{
					Name: "空リスト => 何もしない",
					When: when{
						echos: []model.Echo{},
					},
					Then: func(t *testing.T, dtos []*dao.Echo, err error) {
						require.NoError(t, err)

						want := []*dao.Echo{
							{
								ID:        dtos[0].ID,
								Message:   "m1",
								Timestamp: now,
							},
						}
						if diff := cmp.Diff(want, dtos, cmpopts.IgnoreFields(dao.Echo{}, "CreatedAt", "UpdatedAt")); diff != "" {
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

			var dtos []*dao.Echo
			err := conn.BeginRwTransaction(context.Background(), func(ctx context.Context, tx transaction.Transaction) error {
				r := NewRepository()
				err := r.Save(ctx, tx, when.echos...)
				if err != nil {
					return err
				}

				dtos, err = dao.Echos().All(ctx, tx)
				if err != nil {
					return err
				}
				return nil
			})
			then(t, dtos, err)
		})
	}
}
