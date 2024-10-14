package janken_session_repoimpl_test

import (
	"context"
	"testing"

	"github.com/averak/gamebox/app/adapter/dao"
	"github.com/averak/gamebox/app/adapter/repoimpl/janken_session_repoimpl"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/testutils"
	"github.com/averak/gamebox/testutils/faker"
	"github.com/averak/gamebox/testutils/fixture"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Get(t *testing.T) {
	conn := testutils.MustDBConn(t)

	type args struct {
		gameSessionID uuid.UUID
	}
	tests := []struct {
		name    string
		seeds   []fixture.Seed
		args    args
		want    model.JankenSession
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "レコードが存在する => 取得できる",
			seeds: []fixture.Seed{
				&dao.User{
					ID: faker.UUIDv5("u1").String(),
				},
				&dao.UserGameSession{
					ID:     faker.UUIDv5("gs1").String(),
					UserID: faker.UUIDv5("u1").String(),
					GameID: int(model.GameIDJanken),
				},
				&dao.UserJankenSession{
					GameSessionID: faker.UUIDv5("gs1").String(),
					Seed:          100,
				},
				&dao.UserJankenSessionHistory{
					GameSessionID: faker.UUIDv5("gs1").String(),
					Turn:          1,
					MyHand:        int(model.JankenHandRock),
					OpponentHand:  int(model.JankenHandRock),
				},
				&dao.UserJankenSessionHistory{
					GameSessionID: faker.UUIDv5("gs1").String(),
					Turn:          2,
					MyHand:        int(model.JankenHandScissors),
					OpponentHand:  int(model.JankenHandPaper),
				},
			},
			args: args{
				gameSessionID: faker.UUIDv5("gs1"),
			},
			want: model.JankenSession{
				GameSessionID: faker.UUIDv5("gs1"),
				Seed:          100,
				Histories: []model.JankenHistory{
					{
						Turn:         1,
						MyHand:       model.JankenHandRock,
						OpponentHand: model.JankenHandRock,
					},
					{
						Turn:         2,
						MyHand:       model.JankenHandScissors,
						OpponentHand: model.JankenHandPaper,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:  "レコードが存在しない => エラー",
			seeds: []fixture.Seed{},
			args: args{
				gameSessionID: faker.UUIDv5("gs1"),
			},
			want: model.JankenSession{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, repository.ErrJankenSessionNotFound)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture.SetupSeeds(t, context.Background(), tt.seeds...)
			defer testutils.Teardown(t)

			var got model.JankenSession
			err := conn.BeginRoTransaction(context.Background(), func(ctx context.Context, tx transaction.Transaction) error {
				r := janken_session_repoimpl.NewRepository()
				var err error
				got, err = r.Get(ctx, tx, tt.args.gameSessionID)
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
		session model.JankenSession
	}
	tests := []struct {
		name    string
		seeds   []fixture.Seed
		args    args
		want    model.JankenSession
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "レコードが存在しない => 作成できる",
			seeds: []fixture.Seed{
				&dao.User{
					ID: faker.UUIDv5("u1").String(),
				},
				&dao.UserGameSession{
					ID:     faker.UUIDv5("gs1").String(),
					UserID: faker.UUIDv5("u1").String(),
					GameID: int(model.GameIDJanken),
				},
			},
			args: args{
				session: model.JankenSession{
					GameSessionID: faker.UUIDv5("gs1"),
					Seed:          100,
					Histories: []model.JankenHistory{
						{
							Turn:         1,
							MyHand:       model.JankenHandRock,
							OpponentHand: model.JankenHandRock,
						},
						{
							Turn:         2,
							MyHand:       model.JankenHandScissors,
							OpponentHand: model.JankenHandPaper,
						},
					},
				},
			},
			want: model.JankenSession{
				GameSessionID: faker.UUIDv5("gs1"),
				Seed:          100,
				Histories: []model.JankenHistory{
					{
						Turn:         1,
						MyHand:       model.JankenHandRock,
						OpponentHand: model.JankenHandRock,
					},
					{
						Turn:         2,
						MyHand:       model.JankenHandScissors,
						OpponentHand: model.JankenHandPaper,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "レコードが存在する => 更新できる",
			seeds: []fixture.Seed{
				&dao.User{
					ID: faker.UUIDv5("u1").String(),
				},
				&dao.UserGameSession{
					ID:     faker.UUIDv5("gs1").String(),
					UserID: faker.UUIDv5("u1").String(),
					GameID: int(model.GameIDJanken),
				},
				&dao.UserJankenSession{
					GameSessionID: faker.UUIDv5("gs1").String(),
					Seed:          100,
				},
				&dao.UserJankenSessionHistory{
					GameSessionID: faker.UUIDv5("gs1").String(),
					Turn:          1,
					MyHand:        int(model.JankenHandRock),
					OpponentHand:  int(model.JankenHandRock),
				},
				&dao.UserJankenSessionHistory{
					GameSessionID: faker.UUIDv5("gs1").String(),
					Turn:          2,
					MyHand:        int(model.JankenHandScissors),
					OpponentHand:  int(model.JankenHandPaper),
				},
			},
			args: args{
				session: model.JankenSession{
					GameSessionID: faker.UUIDv5("gs1"),
					Seed:          200,
					// 実際には既存の履歴が書き換わることはないが、テストのために書き換えている。
					Histories: []model.JankenHistory{
						{
							Turn:         1,
							MyHand:       model.JankenHandScissors,
							OpponentHand: model.JankenHandScissors,
						},
						{
							Turn:         3,
							MyHand:       model.JankenHandScissors,
							OpponentHand: model.JankenHandPaper,
						},
					},
				},
			},
			want: model.JankenSession{
				GameSessionID: faker.UUIDv5("gs1"),
				Seed:          200,
				Histories: []model.JankenHistory{
					{
						// 実際には既存の履歴が書き換わることはないが、テストのために書き換えている。
						Turn:         1,
						MyHand:       model.JankenHandScissors,
						OpponentHand: model.JankenHandScissors,
					},
					{
						Turn:         3,
						MyHand:       model.JankenHandScissors,
						OpponentHand: model.JankenHandPaper,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture.SetupSeeds(t, context.Background(), tt.seeds...)
			defer testutils.Teardown(t)

			var got model.JankenSession
			err := conn.BeginRwTransaction(context.Background(), func(ctx context.Context, tx transaction.Transaction) error {
				r := janken_session_repoimpl.NewRepository()
				err := r.Save(ctx, tx, tt.args.session)
				if err != nil {
					return err
				}

				got, err = r.Get(ctx, tx, tt.args.session.GameSessionID)
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

func Test_CheckHistoryDiff(t *testing.T) {
	type args struct {
		newDtos     []*dao.UserJankenSessionHistory
		currentDtos []*dao.UserJankenSessionHistory
	}
	tests := []struct {
		name         string
		args         args
		wantUpserted []*dao.UserJankenSessionHistory
		wantDeleted  []*dao.UserJankenSessionHistory
	}{
		{
			name: "作成/更新/削除された履歴が存在する => upserted と deleted が正しく計算される",
			args: args{
				// gs1: 更新
				// gs2: 削除
				// gs3: 作成
				newDtos: []*dao.UserJankenSessionHistory{
					{
						GameSessionID: faker.UUIDv5("gs1").String(),
						Turn:          1,
						MyHand:        int(model.JankenHandScissors),
					},
					{
						GameSessionID: faker.UUIDv5("gs1").String(),
						Turn:          3,
						MyHand:        int(model.JankenHandRock),
					},
				},
				currentDtos: []*dao.UserJankenSessionHistory{
					{
						GameSessionID: faker.UUIDv5("gs1").String(),
						Turn:          1,
						MyHand:        int(model.JankenHandRock),
					},
					{
						GameSessionID: faker.UUIDv5("gs1").String(),
						Turn:          2,
						MyHand:        int(model.JankenHandRock),
					},
				},
			},
			wantUpserted: []*dao.UserJankenSessionHistory{
				{
					GameSessionID: faker.UUIDv5("gs1").String(),
					Turn:          1,
					MyHand:        int(model.JankenHandScissors),
				},
				{
					GameSessionID: faker.UUIDv5("gs1").String(),
					Turn:          3,
					MyHand:        int(model.JankenHandRock),
				},
			},
			wantDeleted: []*dao.UserJankenSessionHistory{
				{
					GameSessionID: faker.UUIDv5("gs1").String(),
					Turn:          2,
					MyHand:        int(model.JankenHandRock),
				},
			},
		},
		{
			name: "全ての履歴が完全一致 => 何もしない",
			args: args{
				newDtos: []*dao.UserJankenSessionHistory{
					{
						GameSessionID: faker.UUIDv5("gs1").String(),
						Turn:          1,
						MyHand:        int(model.JankenHandRock),
					},
					{
						GameSessionID: faker.UUIDv5("gs1").String(),
						Turn:          2,
						MyHand:        int(model.JankenHandRock),
					},
				},
				currentDtos: []*dao.UserJankenSessionHistory{
					{
						GameSessionID: faker.UUIDv5("gs1").String(),
						Turn:          1,
						MyHand:        int(model.JankenHandRock),
					},
					{
						GameSessionID: faker.UUIDv5("gs1").String(),
						Turn:          2,
						MyHand:        int(model.JankenHandRock),
					},
				},
			},
			wantUpserted: []*dao.UserJankenSessionHistory{},
			wantDeleted:  []*dao.UserJankenSessionHistory{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUpserted, gotDeleted := janken_session_repoimpl.CheckHistoryDiff(tt.args.newDtos, tt.args.currentDtos)
			assert.Equal(t, tt.wantUpserted, gotUpserted)
			assert.Equal(t, tt.wantDeleted, gotDeleted)
		})
	}
}
