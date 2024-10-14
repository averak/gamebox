package handler

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/app/adapter/handler/debug/echo_handler"
	"github.com/averak/gamebox/app/adapter/handler/game_handler"
	"github.com/averak/gamebox/app/adapter/handler/janken_handler"
	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/infrastructure/connect/interceptor"
	"github.com/averak/gamebox/protobuf/api/apiconnect"
	"github.com/averak/gamebox/protobuf/api/debug/debugconnect"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	game_handler.NewHandler,
	janken_handler.NewHandler,
	echo_handler.NewHandler,
	New,
)

func New(
	gameHandler apiconnect.GameServiceHandler,
	jankenHandler apiconnect.JankenServiceHandler,
	echoHandler debugconnect.EchoServiceHandler,
) *http.ServeMux {
	opts := connect.WithInterceptors(interceptor.New()...)
	mux := http.NewServeMux()
	mux.Handle(apiconnect.NewGameServiceHandler(gameHandler, opts))
	mux.Handle(apiconnect.NewJankenServiceHandler(jankenHandler, opts))
	if config.Get().GetDebug() {
		mux.Handle(debugconnect.NewEchoServiceHandler(echoHandler, opts))
	}
	return mux
}
