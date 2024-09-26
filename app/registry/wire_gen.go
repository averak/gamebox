// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package registry

import (
	"context"
	"github.com/averak/gamebox/app/adapter/handler"
	"net/http"
)

// Injectors from wire.go:

func InitializeAPIServerMux(ctx context.Context) (*http.ServeMux, error) {
	serveMux := handler.New()
	return serveMux, nil
}
