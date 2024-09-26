package handler

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/app/adapter/handler/debug/echo_handler"
	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/infrastructure/connect/interceptor"
	"github.com/averak/gamebox/protobuf/api/debug/debugconnect"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	echo_handler.NewHandler,
	New,
)

func New(echo debugconnect.EchoServiceHandler) *http.ServeMux {
	opts := connect.WithInterceptors(interceptor.New()...)
	mux := http.NewServeMux()
	if config.Get().GetDebug() {
		mux.Handle(debugconnect.NewEchoServiceHandler(echo, opts))
	}
	return mux
}
