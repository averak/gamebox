package user_repoimpl

import (
	"context"
	"testing"

	"github.com/averak/gamebox/app/adapter/dao"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/testutils"
	"github.com/averak/gamebox/testutils/faker"
	"github.com/averak/gamebox/testutils/fixture"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Get(t *testing.T) {
	conn := testutils.MustDBConn(t)

	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		seeds   []fixture.Seed
		args    args
		want    model.User
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "レコードが存在する => 取得できる",
			seeds: []fixture.Seed{
				&dao.User{
					ID: faker.UUIDv5("u1").String(),
				},
			},
			args: args{
				userID: faker.UUIDv5("u1"),
			},
			want: model.User{
				ID: faker.UUIDv5("u1"),
			},
			wantErr: assert.NoError,
		},
		{
			name:  "レコードが存在しない => エラー",
			seeds: []fixture.Seed{},
			args: args{
				userID: faker.UUIDv5("not exists"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, repository.ErrUserNotFound)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture.SetupSeeds(t, context.Background(), tt.seeds...)
			defer testutils.Teardown(t)

			var got model.User
			err := conn.BeginRoTransaction(context.Background(), func(ctx context.Context, tx transaction.Transaction) error {
				r := Repository{}
				var err error
				got, err = r.Get(ctx, tx, tt.args.userID)
				if err != nil {
					return err
				}
				return nil
			})
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRepository_Save(t *testing.T) {
	conn := testutils.MustDBConn(t)

	type args struct {
		user model.User
	}
	tests := []struct {
		name    string
		seeds   []fixture.Seed
		args    args
		want    []*dao.User
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:  "PK が存在しない => 作成する",
			seeds: []fixture.Seed{},
			args: args{
				user: model.User{
					ID:     faker.UUIDv5("u1"),
					Status: model.UserStatusActive,
				},
			},
			want: []*dao.User{
				{
					ID:     faker.UUIDv5("u1").String(),
					Status: int(model.UserStatusActive),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "PK が存在する => 更新する",
			seeds: []fixture.Seed{
				&dao.User{
					ID:     faker.UUIDv5("u1").String(),
					Status: int(model.UserStatusActive),
				},
			},
			args: args{
				user: model.User{
					ID:     faker.UUIDv5("u1"),
					Status: model.UserStatusDeactivated,
				},
			},
			want: []*dao.User{
				{
					ID:     faker.UUIDv5("u1").String(),
					Status: int(model.UserStatusDeactivated),
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture.SetupSeeds(t, context.Background(), tt.seeds...)
			defer testutils.Teardown(t)

			var got []*dao.User
			err := conn.BeginRwTransaction(context.Background(), func(ctx context.Context, tx transaction.Transaction) error {
				r := Repository{}
				err := r.Save(ctx, tx, tt.args.user)
				if err != nil {
					return err
				}

				got, err = dao.Users().All(ctx, tx)
				if err != nil {
					return err
				}
				return nil
			})
			if !tt.wantErr(t, err) {
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(dao.User{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}
