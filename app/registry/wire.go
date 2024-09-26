//go:build wireinject
// +build wireinject

package registry

import (
	"context"
	"net/http"

	"github.com/averak/gamebox/app/adapter/handler"
	"github.com/google/wire"
)

func InitializeAPIServerMux(ctx context.Context) (*http.ServeMux, error) {
	wire.Build(handler.SuperSet)
	return nil, nil
}
