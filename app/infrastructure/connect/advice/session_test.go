package advice

import (
	"context"
	"net/http"
	"net/textproto"
	"testing"
	"time"

	"github.com/averak/gamebox/app/adapter/dao"
	"github.com/averak/gamebox/app/adapter/repoimpl/user_repoimpl"
	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/app/infrastructure/connect/mdval"
	"github.com/averak/gamebox/app/infrastructure/session"
	"github.com/averak/gamebox/protobuf/api/api_errors"
	pb "github.com/averak/gamebox/protobuf/config"
	"github.com/averak/gamebox/testutils"
	"github.com/averak/gamebox/testutils/bdd"
	"github.com/averak/gamebox/testutils/faker"
	"github.com/averak/gamebox/testutils/fixture"
	"github.com/averak/gamebox/testutils/testconnect"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_checkSession(t *testing.T) {
	now := time.Now()

	sessionToken := func(userID uuid.UUID) string {
		sessionToken, err := session.EncodeSessionToken(
			session.NewSession(userID, now, now.Add(1*time.Hour)),
			[]byte(config.Get().GetApiServer().GetSession().GetSecretKey()),
		)
		if err != nil {
			t.Fatal(err)
		}
		return sessionToken
	}

	type given struct {
		conf  *config.Config
		seeds []fixture.Seed
	}
	type when struct {
		incomingMD mdval.IncomingMD
		now        time.Time
	}
	type then = func(*testing.T, *model.User, error)
	tests := []bdd.Testcase[given, when, then]{
		{
			Name: "デバッグモードの場合",
			Given: given{
				conf: &config.Config{
					Debug: true,
				},
				seeds: []fixture.Seed{
					&dao.User{
						ID:     faker.UUIDv5("u1").String(),
						Status: int(model.UserStatusActive),
					},
				},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "デバッグ用ヘッダが指定されている && ユーザが存在する場合 => ユーザ情報を取得できる",
					When: when{
						incomingMD: mdval.NewIncomingMD(http.Header{
							textproto.CanonicalMIMEHeaderKey(string(mdval.DebugSpoofingUserIDKey)): {faker.UUIDv5("u1").String()},
						}),
						now: now,
					},
					Then: func(t *testing.T, got *model.User, err error) {
						require.NoError(t, err)

						want := &model.User{
							ID:     faker.UUIDv5("u1"),
							Status: model.UserStatusActive,
						}
						assert.Equal(t, want, got)
					},
				},
				{
					Name: "デバッグ用ヘッダが指定されている && ユーザが存在しない場合 => ユーザを新規作成できる",
					When: when{
						incomingMD: mdval.NewIncomingMD(http.Header{
							textproto.CanonicalMIMEHeaderKey(string(mdval.DebugSpoofingUserIDKey)): {faker.UUIDv5("u2").String()},
						}),
						now: now,
					},
					Then: func(t *testing.T, got *model.User, err error) {
						require.NoError(t, err)

						want := &model.User{
							ID:     faker.UUIDv5("u2"),
							Status: model.UserStatusActive,
						}
						assert.Equal(t, want, got)
					},
				},
			},
		},
		{
			Name: "デバッグモードではない場合",
			Given: given{
				conf: &config.Config{
					Debug: false,
					ApiServer: &pb.APIServer{
						Session: &pb.APIServer_Session{
							SecretKey: config.Get().GetApiServer().GetSession().GetSecretKey(),
						},
					},
				},
				seeds: []fixture.Seed{
					&dao.User{
						ID: faker.UUIDv5("u1").String(),
					},
					&dao.User{
						ID: faker.UUIDv5("u2").String(),
					},
				},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "セッショントークンが指定されている && ユーザが存在する場合 => ユーザ情報を取得できる",
					When: when{
						incomingMD: mdval.NewIncomingMD(http.Header{
							textproto.CanonicalMIMEHeaderKey(string(mdval.SessionTokenKey)): {sessionToken(faker.UUIDv5("u1"))},
						}),
						now: now,
					},
					Then: func(t *testing.T, got *model.User, err error) {
						require.NoError(t, err)

						want := &model.User{
							ID: faker.UUIDv5("u1"),
						}
						assert.Equal(t, want, got)
					},
				},
				{
					Name: "セッショントークンが指定されている && ユーザが存在しない場合 => エラー",
					When: when{
						incomingMD: mdval.NewIncomingMD(http.Header{
							textproto.CanonicalMIMEHeaderKey(string(mdval.SessionTokenKey)): {sessionToken(faker.UUIDv5("not_exists"))},
						}),
						now: now,
					},
					Then: func(t *testing.T, got *model.User, err error) {
						testconnect.AssertErrorCode(t, api_errors.ErrorCode_COMMON_INVALID_SESSION, err)
					},
				},
				{
					Name: "セッショントークンが有効期限切れ => エラー",
					When: when{
						incomingMD: mdval.NewIncomingMD(http.Header{
							textproto.CanonicalMIMEHeaderKey(string(mdval.SessionTokenKey)): {sessionToken(faker.UUIDv5("u1"))},
						}),
						now: now.Add(1 * time.Hour),
					},
					Then: func(t *testing.T, got *model.User, err error) {
						testconnect.AssertErrorCode(t, api_errors.ErrorCode_COMMON_INVALID_SESSION, err)
					},
				},
				{
					Name: "不正なセッショントークン => エラー",
					When: when{
						incomingMD: mdval.NewIncomingMD(http.Header{
							textproto.CanonicalMIMEHeaderKey(string(mdval.SessionTokenKey)): {"invalid"},
						}),
						now: now,
					},
					Then: func(t *testing.T, got *model.User, err error) {
						testconnect.AssertErrorCode(t, api_errors.ErrorCode_COMMON_INVALID_SESSION, err)
					},
				},
				{
					Name: "デバッグ用ヘッダは無視される",
					When: when{
						incomingMD: mdval.NewIncomingMD(http.Header{
							textproto.CanonicalMIMEHeaderKey(string(mdval.DebugSpoofingUserIDKey)): {faker.UUIDv5("u1").String()},
						}),
						now: now,
					},
					Then: func(t *testing.T, got *model.User, err error) {
						testconnect.AssertErrorCode(t, api_errors.ErrorCode_COMMON_INVALID_SESSION, err)
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.Run(t, func(t *testing.T, given given, when when, then then) {
			fixture.SetupSeeds(t, context.Background(), given.seeds...)
			defer testutils.Teardown(t)

			conn := testutils.MustDBConn(t)

			got, err := checkSession(context.Background(), given.conf, user_repoimpl.NewRepository(), conn, when.incomingMD, when.now)
			then(t, got, err)
		})
	}
}

func Test_setupSpoofingUser(t *testing.T) {
	conn := testutils.MustDBConn(t)

	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name     string
		seeds    []fixture.Seed
		args     args
		want     model.User
		wantDtos []*dao.User
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "ユーザが存在する場合 => ユーザ情報を取得できる",
			seeds: []fixture.Seed{
				&dao.User{
					ID:     faker.UUIDv5("u1").String(),
					Status: int(model.UserStatusActive),
				},
			},
			args: args{
				userID: faker.UUIDv5("u1"),
			},
			want: model.User{
				ID:     faker.UUIDv5("u1"),
				Status: model.UserStatusActive,
			},
			wantDtos: []*dao.User{
				{
					ID:     faker.UUIDv5("u1").String(),
					Status: int(model.UserStatusActive),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:  "ユーザが存在しない場合 => ユーザを作成できる",
			seeds: []fixture.Seed{},
			args: args{
				userID: faker.UUIDv5("u1"),
			},
			want: model.User{
				ID:     faker.UUIDv5("u1"),
				Status: model.UserStatusActive,
			},
			wantDtos: []*dao.User{
				{
					ID:     faker.UUIDv5("u1").String(),
					Status: int(model.UserStatusActive),
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture.SetupSeeds(t, context.Background(), tt.seeds...)
			defer testutils.Teardown(t)

			got, err := setupSpoofingUser(context.Background(), conn, user_repoimpl.NewRepository(), tt.args.userID)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)

			var dtos []*dao.User
			eerr := conn.BeginRoTransaction(context.Background(), func(ctx context.Context, tx transaction.Transaction) error {
				var err error
				dtos, err = dao.Users().All(ctx, tx)
				return err
			})
			if eerr != nil {
				t.Fatal(eerr)
			}
			if diff := cmp.Diff(tt.wantDtos, dtos, cmpopts.IgnoreFields(dao.User{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}
