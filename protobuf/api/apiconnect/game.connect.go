// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: api/game.proto

package apiconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	api "github.com/averak/gamebox/protobuf/api"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// GameServiceName is the fully-qualified name of the GameService service.
	GameServiceName = "api.GameService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// GameServiceGetSessionV1Procedure is the fully-qualified name of the GameService's GetSessionV1
	// RPC.
	GameServiceGetSessionV1Procedure = "/api.GameService/GetSessionV1"
	// GameServiceListPlayingSessionsV1Procedure is the fully-qualified name of the GameService's
	// ListPlayingSessionsV1 RPC.
	GameServiceListPlayingSessionsV1Procedure = "/api.GameService/ListPlayingSessionsV1"
	// GameServiceStartPlayingV1Procedure is the fully-qualified name of the GameService's
	// StartPlayingV1 RPC.
	GameServiceStartPlayingV1Procedure = "/api.GameService/StartPlayingV1"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	gameServiceServiceDescriptor                     = api.File_api_game_proto.Services().ByName("GameService")
	gameServiceGetSessionV1MethodDescriptor          = gameServiceServiceDescriptor.Methods().ByName("GetSessionV1")
	gameServiceListPlayingSessionsV1MethodDescriptor = gameServiceServiceDescriptor.Methods().ByName("ListPlayingSessionsV1")
	gameServiceStartPlayingV1MethodDescriptor        = gameServiceServiceDescriptor.Methods().ByName("StartPlayingV1")
)

// GameServiceClient is a client for the api.GameService service.
type GameServiceClient interface {
	// ゲームセッションを取得します。
	GetSessionV1(context.Context, *connect.Request[api.GameServiceGetSessionV1Request]) (*connect.Response[api.GameServiceGetSessionV1Response], error)
	// プレイ中のゲームセッションリストを取得します。
	ListPlayingSessionsV1(context.Context, *connect.Request[api.GameServiceListPlayingSessionsV1Request]) (*connect.Response[api.GameServiceListPlayingSessionsV1Response], error)
	// ゲームを開始します。
	// なお、チート対策のためゲーム終了判定はサーバー側で行います。
	StartPlayingV1(context.Context, *connect.Request[api.GameServiceStartPlayingV1Request]) (*connect.Response[api.GameServiceStartPlayingV1Response], error)
}

// NewGameServiceClient constructs a client for the api.GameService service. By default, it uses the
// Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewGameServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) GameServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &gameServiceClient{
		getSessionV1: connect.NewClient[api.GameServiceGetSessionV1Request, api.GameServiceGetSessionV1Response](
			httpClient,
			baseURL+GameServiceGetSessionV1Procedure,
			connect.WithSchema(gameServiceGetSessionV1MethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listPlayingSessionsV1: connect.NewClient[api.GameServiceListPlayingSessionsV1Request, api.GameServiceListPlayingSessionsV1Response](
			httpClient,
			baseURL+GameServiceListPlayingSessionsV1Procedure,
			connect.WithSchema(gameServiceListPlayingSessionsV1MethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		startPlayingV1: connect.NewClient[api.GameServiceStartPlayingV1Request, api.GameServiceStartPlayingV1Response](
			httpClient,
			baseURL+GameServiceStartPlayingV1Procedure,
			connect.WithSchema(gameServiceStartPlayingV1MethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// gameServiceClient implements GameServiceClient.
type gameServiceClient struct {
	getSessionV1          *connect.Client[api.GameServiceGetSessionV1Request, api.GameServiceGetSessionV1Response]
	listPlayingSessionsV1 *connect.Client[api.GameServiceListPlayingSessionsV1Request, api.GameServiceListPlayingSessionsV1Response]
	startPlayingV1        *connect.Client[api.GameServiceStartPlayingV1Request, api.GameServiceStartPlayingV1Response]
}

// GetSessionV1 calls api.GameService.GetSessionV1.
func (c *gameServiceClient) GetSessionV1(ctx context.Context, req *connect.Request[api.GameServiceGetSessionV1Request]) (*connect.Response[api.GameServiceGetSessionV1Response], error) {
	return c.getSessionV1.CallUnary(ctx, req)
}

// ListPlayingSessionsV1 calls api.GameService.ListPlayingSessionsV1.
func (c *gameServiceClient) ListPlayingSessionsV1(ctx context.Context, req *connect.Request[api.GameServiceListPlayingSessionsV1Request]) (*connect.Response[api.GameServiceListPlayingSessionsV1Response], error) {
	return c.listPlayingSessionsV1.CallUnary(ctx, req)
}

// StartPlayingV1 calls api.GameService.StartPlayingV1.
func (c *gameServiceClient) StartPlayingV1(ctx context.Context, req *connect.Request[api.GameServiceStartPlayingV1Request]) (*connect.Response[api.GameServiceStartPlayingV1Response], error) {
	return c.startPlayingV1.CallUnary(ctx, req)
}

// GameServiceHandler is an implementation of the api.GameService service.
type GameServiceHandler interface {
	// ゲームセッションを取得します。
	GetSessionV1(context.Context, *connect.Request[api.GameServiceGetSessionV1Request]) (*connect.Response[api.GameServiceGetSessionV1Response], error)
	// プレイ中のゲームセッションリストを取得します。
	ListPlayingSessionsV1(context.Context, *connect.Request[api.GameServiceListPlayingSessionsV1Request]) (*connect.Response[api.GameServiceListPlayingSessionsV1Response], error)
	// ゲームを開始します。
	// なお、チート対策のためゲーム終了判定はサーバー側で行います。
	StartPlayingV1(context.Context, *connect.Request[api.GameServiceStartPlayingV1Request]) (*connect.Response[api.GameServiceStartPlayingV1Response], error)
}

// NewGameServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewGameServiceHandler(svc GameServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	gameServiceGetSessionV1Handler := connect.NewUnaryHandler(
		GameServiceGetSessionV1Procedure,
		svc.GetSessionV1,
		connect.WithSchema(gameServiceGetSessionV1MethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gameServiceListPlayingSessionsV1Handler := connect.NewUnaryHandler(
		GameServiceListPlayingSessionsV1Procedure,
		svc.ListPlayingSessionsV1,
		connect.WithSchema(gameServiceListPlayingSessionsV1MethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gameServiceStartPlayingV1Handler := connect.NewUnaryHandler(
		GameServiceStartPlayingV1Procedure,
		svc.StartPlayingV1,
		connect.WithSchema(gameServiceStartPlayingV1MethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/api.GameService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case GameServiceGetSessionV1Procedure:
			gameServiceGetSessionV1Handler.ServeHTTP(w, r)
		case GameServiceListPlayingSessionsV1Procedure:
			gameServiceListPlayingSessionsV1Handler.ServeHTTP(w, r)
		case GameServiceStartPlayingV1Procedure:
			gameServiceStartPlayingV1Handler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedGameServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedGameServiceHandler struct{}

func (UnimplementedGameServiceHandler) GetSessionV1(context.Context, *connect.Request[api.GameServiceGetSessionV1Request]) (*connect.Response[api.GameServiceGetSessionV1Response], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.GameService.GetSessionV1 is not implemented"))
}

func (UnimplementedGameServiceHandler) ListPlayingSessionsV1(context.Context, *connect.Request[api.GameServiceListPlayingSessionsV1Request]) (*connect.Response[api.GameServiceListPlayingSessionsV1Response], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.GameService.ListPlayingSessionsV1 is not implemented"))
}

func (UnimplementedGameServiceHandler) StartPlayingV1(context.Context, *connect.Request[api.GameServiceStartPlayingV1Request]) (*connect.Response[api.GameServiceStartPlayingV1Response], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.GameService.StartPlayingV1 is not implemented"))
}
