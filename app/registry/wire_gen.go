// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package registry

import (
	"context"
	"github.com/averak/gamebox/app/adapter/handler"
	"github.com/averak/gamebox/app/adapter/handler/debug/echo_handler"
	"github.com/averak/gamebox/app/adapter/repoimpl"
	"github.com/averak/gamebox/app/adapter/repoimpl/echo_repoimpl"
	"github.com/averak/gamebox/app/infrastructure/connect/advice"
	"github.com/averak/gamebox/app/infrastructure/db"
	"github.com/averak/gamebox/app/usecase"
	"github.com/averak/gamebox/app/usecase/echo_usecase"
	"github.com/google/wire"
	"net/http"
)

// Injectors from wire.go:

func InitializeAPIServerMux(ctx context.Context) (*http.ServeMux, error) {
	connection, err := db.NewConnection()
	if err != nil {
		return nil, err
	}
	echoRepository := echo_repoimpl.NewRepository()
	usecase := echo_usecase.NewUsecase(connection, echoRepository)
	adviceAdvice := advice.NewAdvice(connection)
	echoServiceHandler := echo_handler.NewHandler(usecase, adviceAdvice)
	serveMux := handler.New(echoServiceHandler)
	return serveMux, nil
}

// wire.go:

var SuperSet = wire.NewSet(repoimpl.SuperSet, usecase.SuperSet, advice.NewAdvice, db.NewConnection)
