//go:build wireinject
// +build wireinject

package registry

import (
	"context"
	"net/http"

	"github.com/averak/gamebox/app/adapter/handler"
	"github.com/averak/gamebox/app/adapter/repoimpl"
	"github.com/averak/gamebox/app/infrastructure/connect/advice"
	"github.com/averak/gamebox/app/infrastructure/db"
	"github.com/averak/gamebox/app/usecase"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	repoimpl.SuperSet,
	usecase.SuperSet,
	advice.NewAdvice,
	db.NewConnection,
)

func InitializeAPIServerMux(ctx context.Context) (*http.ServeMux, error) {
	wire.Build(SuperSet, handler.SuperSet)
	return nil, nil
}
