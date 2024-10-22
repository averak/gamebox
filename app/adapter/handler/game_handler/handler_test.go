package game_handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/registry"
	"github.com/averak/gamebox/protobuf/api"
	"github.com/averak/gamebox/protobuf/api/api_errors"
	"github.com/averak/gamebox/protobuf/api/apiconnect"
	"github.com/averak/gamebox/protobuf/resource"
	"github.com/averak/gamebox/testutils"
	"github.com/averak/gamebox/testutils/bdd"
	"github.com/averak/gamebox/testutils/faker"
	"github.com/averak/gamebox/testutils/fixture"
	"github.com/averak/gamebox/testutils/fixture/fixture_builder"
	"github.com/averak/gamebox/testutils/testconnect"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_handler_GetSessionV1(t *testing.T) {
	mux, err := registry.InitializeAPIServerMux(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	now := time.Now().Truncate(time.Millisecond)

	type given struct {
		fixtures []fixture_builder.Fixture
	}
	type when struct {
		req  *api.GameServiceGetSessionV1Request
		opts []testconnect.Option
	}
	type then = func(*testing.T, *connect.Response[api.GameServiceGetSessionV1Response], error)
	tests := []bdd.Testcase[given, when, then]{
		{
			Name: "ゲームセッションが存在する状態で",
			Given: given{
				fixtures: []fixture_builder.Fixture{
					fixture_builder.New().
						User(fixture_builder.NewUserBuilder(faker.UUIDv5("u1")).Build()).
						GameSession(
							fixture_builder.NewGameSessionBuilder(t, faker.UUIDv5("u1"), faker.UUIDv5("u1_gs1")).
								StartedAt(now.Add(-time.Hour)).
								FinishedAt(now).
								Build(),
							fixture_builder.NewGameSessionBuilder(t, faker.UUIDv5("u1"), faker.UUIDv5("u1_gs2")).Build(),
						).
						Build(),
					fixture_builder.New().
						User(fixture_builder.NewUserBuilder(faker.UUIDv5("u2")).Build()).
						GameSession(
							fixture_builder.NewGameSessionBuilder(t, faker.UUIDv5("u2"), faker.UUIDv5("u2_gs1")).Build(),
						).
						Build(),
				},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "[正常系] 自分のゲームセッションIDを指定 => ゲームセッション情報を取得できる",
					When: when{
						req: &api.GameServiceGetSessionV1Request{
							SessionId: faker.UUIDv5("u1_gs1").String(),
						},
						opts: []testconnect.Option{
							testconnect.WithSpoofingUserID(faker.UUIDv5("u1")),
						},
					},
					Then: func(t *testing.T, got *connect.Response[api.GameServiceGetSessionV1Response], err error) {
						require.NoError(t, err)

						want := &api.GameServiceGetSessionV1Response{
							Session: &resource.GameSession{
								SessionId:  faker.UUIDv5("u1_gs1").String(),
								GameId:     resource.GameID_GAME_ID_SOLITAIRE,
								Status:     resource.GameStatus_GAME_STATUS_FINISHED,
								Result:     resource.GameResult_GAME_RESULT_WIN,
								Wager:      100,
								Payout:     200,
								StartedAt:  timestamppb.New(now.Add(-time.Hour)),
								FinishedAt: timestamppb.New(now),
							},
						}
						assert.EqualExportedValues(t, want, got.Msg)
					},
				},
				{
					Name: "[異常系] 他人のゲームセッションIDを指定 => エラー",
					When: when{
						req: &api.GameServiceGetSessionV1Request{
							SessionId: faker.UUIDv5("u2_gs1").String(),
						},
						opts: []testconnect.Option{
							testconnect.WithSpoofingUserID(faker.UUIDv5("u1")),
						},
					},
					Then: func(t *testing.T, got *connect.Response[api.GameServiceGetSessionV1Response], err error) {
						testconnect.AssertErrorCode(t, api_errors.ErrorCode_METHOD_RESOURCE_NOT_FOUND, err)
					},
				},
				{
					Name: "[異常系] 存在しないIDを指定 => エラー",
					When: when{
						req: &api.GameServiceGetSessionV1Request{
							SessionId: faker.UUIDv5("not exist").String(),
						},
						opts: []testconnect.Option{
							testconnect.WithSpoofingUserID(faker.UUIDv5("u1")),
						},
					},
					Then: func(t *testing.T, got *connect.Response[api.GameServiceGetSessionV1Response], err error) {
						testconnect.AssertErrorCode(t, api_errors.ErrorCode_METHOD_RESOURCE_NOT_FOUND, err)
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.Run(t, func(t *testing.T, given given, when when, then then) {
			fixture.Setup(t, given.fixtures...)
			defer testutils.Teardown(t)

			got, err := testconnect.MethodInvoke(
				apiconnect.NewGameServiceClient(http.DefaultClient, server.URL).GetSessionV1,
				when.req,
				when.opts...,
			)
			then(t, got, err)
		})
	}
}

func Test_handler_ListPlayingSessionsV1(t *testing.T) {
	mux, err := registry.InitializeAPIServerMux(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	type given struct {
		fixtures []fixture_builder.Fixture
	}
	type when struct {
		req  *api.GameServiceListPlayingSessionsV1Request
		opts []testconnect.Option
	}
	type then = func(*testing.T, *connect.Response[api.GameServiceListPlayingSessionsV1Response], error)
	tests := []bdd.Testcase[given, when, then]{
		{
			Name: "プレイ中のゲームセッションが存在する状態で",
			Given: given{
				fixtures: []fixture_builder.Fixture{
					fixture_builder.New().
						User(fixture_builder.NewUserBuilder(faker.UUIDv5("u1")).Build()).
						GameSession(
							fixture_builder.NewGameSessionBuilder(t, faker.UUIDv5("u1"), faker.UUIDv5("u1_gs1")).
								GameID(model.GameIDSolitaire).
								Status(model.GameStatusPlaying).
								Build(),
							fixture_builder.NewGameSessionBuilder(t, faker.UUIDv5("u1"), faker.UUIDv5("u1_gs2")).
								GameID(model.GameIDBlackjack).
								Status(model.GameStatusPlaying).
								Build(),
							// 終了したゲームセッションは取得されない
							fixture_builder.NewGameSessionBuilder(t, faker.UUIDv5("u1"), faker.UUIDv5("u1_gs3")).
								GameID(model.GameIDJanken).
								Status(model.GameStatusFinished).
								Build(),
						).
						Build(),
				},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "プレイ中のゲームセッションリストを取得できる",
					When: when{
						req: &api.GameServiceListPlayingSessionsV1Request{},
						opts: []testconnect.Option{
							testconnect.WithSpoofingUserID(faker.UUIDv5("u1")),
						},
					},
					Then: func(t *testing.T, got *connect.Response[api.GameServiceListPlayingSessionsV1Response], err error) {
						require.NoError(t, err)

						want := &api.GameServiceListPlayingSessionsV1Response{
							Sessions: []*resource.GameSession{
								{
									SessionId: faker.UUIDv5("u1_gs1").String(),
									GameId:    resource.GameID_GAME_ID_SOLITAIRE,
									Status:    resource.GameStatus_GAME_STATUS_PLAYING,
								},
								{
									SessionId: faker.UUIDv5("u1_gs2").String(),
									GameId:    resource.GameID_GAME_ID_BLACKJACK,
									Status:    resource.GameStatus_GAME_STATUS_PLAYING,
								},
							},
						}
						if diff := cmp.Diff(want, got.Msg,
							cmpopts.IgnoreUnexported(api.GameServiceListPlayingSessionsV1Response{}, resource.GameSession{}),
							cmpopts.IgnoreFields(resource.GameSession{}, "Result", "Wager", "Payout", "StartedAt", "FinishedAt"),
						); diff != "" {
							t.Errorf("(-want, +got)\n%s", diff)
						}
					},
				},
			},
		},
		{
			Name: "プレイ中のゲームセッションが存在しない状態で",
			Given: given{
				fixtures: []fixture_builder.Fixture{
					fixture_builder.New().
						User(fixture_builder.NewUserBuilder(faker.UUIDv5("u1")).Build()).Build(),
				},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "空リストが返却される",
					When: when{
						req: &api.GameServiceListPlayingSessionsV1Request{},
						opts: []testconnect.Option{
							testconnect.WithSpoofingUserID(faker.UUIDv5("u1")),
						},
					},
					Then: func(t *testing.T, got *connect.Response[api.GameServiceListPlayingSessionsV1Response], err error) {
						require.NoError(t, err)

						want := &api.GameServiceListPlayingSessionsV1Response{}
						assert.EqualExportedValues(t, want, got.Msg)
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.Run(t, func(t *testing.T, given given, when when, then then) {
			fixture.Setup(t, given.fixtures...)
			defer testutils.Teardown(t)

			got, err := testconnect.MethodInvoke(
				apiconnect.NewGameServiceClient(http.DefaultClient, server.URL).ListPlayingSessionsV1,
				when.req,
				when.opts...,
			)
			then(t, got, err)
		})
	}
}

func Test_handler_StartPlayingV1(t *testing.T) {
	mux, err := registry.InitializeAPIServerMux(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	now := time.Now()

	type given struct {
		fixtures []fixture_builder.Fixture
	}
	type when struct {
		req  *api.GameServiceStartPlayingV1Request
		opts []testconnect.Option
	}
	type then = func(*testing.T, *connect.Response[api.GameServiceStartPlayingV1Response], error)
	tests := []bdd.Testcase[given, when, then]{
		{
			Name: "プレイ中/終了済みのゲームセッションが存在する状態で",
			Given: given{
				fixtures: []fixture_builder.Fixture{
					fixture_builder.New().
						User(fixture_builder.NewUserBuilder(faker.UUIDv5("u1")).Build()).
						GameSession(
							fixture_builder.NewGameSessionBuilder(t, faker.UUIDv5("u1"), faker.UUIDv5("u1_gs1")).
								GameID(model.GameIDSolitaire).
								Status(model.GameStatusPlaying).
								Build(),
							fixture_builder.NewGameSessionBuilder(t, faker.UUIDv5("u1"), faker.UUIDv5("u1_gs2")).
								GameID(model.GameIDBlackjack).
								Status(model.GameStatusFinished).
								Build(),
						).
						Build(),
				},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "未プレイ中のゲームIDを指定 => ゲームセッションを開始できる",
					When: when{
						req: &api.GameServiceStartPlayingV1Request{
							GameId: resource.GameID_GAME_ID_BLACKJACK,
							Wager:  100,
						},
						opts: []testconnect.Option{
							testconnect.WithSpoofingUserID(faker.UUIDv5("u1")),
							testconnect.WithIdempotencyKey(faker.UUIDv5("ik1")),
							testconnect.WithAdjustedTime(now),
						},
					},
					Then: func(t *testing.T, got *connect.Response[api.GameServiceStartPlayingV1Response], err error) {
						require.NoError(t, err)

						want := &api.GameServiceStartPlayingV1Response{
							Session: &resource.GameSession{
								SessionId:  faker.UUIDv5("ik1").String(),
								GameId:     resource.GameID_GAME_ID_BLACKJACK,
								Status:     resource.GameStatus_GAME_STATUS_PLAYING,
								Wager:      100,
								Payout:     0,
								StartedAt:  timestamppb.New(now),
								FinishedAt: nil,
							},
						}
						assert.EqualExportedValues(t, want, got.Msg)
					},
				},
				{
					Name: "プレイ中のゲームIDを指定 => エラー",
					When: when{
						req: &api.GameServiceStartPlayingV1Request{
							GameId: resource.GameID_GAME_ID_SOLITAIRE,
							Wager:  100,
						},
						opts: []testconnect.Option{
							testconnect.WithSpoofingUserID(faker.UUIDv5("u1")),
						},
					},
					Then: func(t *testing.T, got *connect.Response[api.GameServiceStartPlayingV1Response], err error) {
						testconnect.AssertErrorCode(t, api_errors.ErrorCode_METHOD_RESOURCE_CONFLICT, err)
					},
				},
				{
					Name: "不正な賭け金 => エラー",
					When: when{
						req: &api.GameServiceStartPlayingV1Request{
							GameId: resource.GameID_GAME_ID_BLACKJACK,
							Wager:  0,
						},
						opts: []testconnect.Option{
							testconnect.WithSpoofingUserID(faker.UUIDv5("u1")),
						},
					},
					Then: func(t *testing.T, got *connect.Response[api.GameServiceStartPlayingV1Response], err error) {
						testconnect.AssertErrorCode(t, api_errors.ErrorCode_METHOD_ILLEGAL_ARGUMENT, err)
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.Run(t, func(t *testing.T, given given, when when, then then) {
			fixture.Setup(t, given.fixtures...)
			defer testutils.Teardown(t)

			got, err := testconnect.MethodInvoke(
				apiconnect.NewGameServiceClient(http.DefaultClient, server.URL).StartPlayingV1,
				when.req,
				when.opts...,
			)
			then(t, got, err)
		})
	}
}
