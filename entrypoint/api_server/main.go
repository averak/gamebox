package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/averak/gamebox/app/core/build_info"
	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/core/logger"
	"github.com/averak/gamebox/app/infrastructure/trace"
	"github.com/averak/gamebox/app/registry"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Init("api-server", build_info.ServerVersion())
	if config.Get().GetGoogleCloud().GetTrace().GetEnabled() {
		trace.Init(ctx, "api-server", build_info.ServerVersion(), config.Get().GetGoogleCloud().GetTrace().GetSamplingRate())
	}

	mux, err := registry.InitializeAPIServerMux(ctx)
	if err != nil {
		logger.Emergency(ctx, map[string]interface{}{
			"message": "Failed to initialize api-server mux.",
			"error":   err.Error(),
		})
		log.Fatal(err)
	}
	mux.Handle(grpchealth.NewHandler(grpchealth.NewStaticChecker()))
	if config.Get().GetDebug() {
		reflector := grpcreflect.NewStaticReflector()
		mux.Handle(grpcreflect.NewHandlerV1(reflector))
		mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	}
	handler := cors.New(cors.Options{
		AllowedOrigins:   config.Get().GetApiServer().GetCors().GetAllowedOrigins(),
		AllowedMethods:   config.Get().GetApiServer().GetCors().GetAllowedMethods(),
		AllowedHeaders:   config.Get().GetApiServer().GetCors().GetAllowedHeaders(),
		ExposedHeaders:   config.Get().GetApiServer().GetCors().GetExposeHeaders(),
		MaxAge:           int(config.Get().GetApiServer().GetCors().GetMaxAge()),
		AllowCredentials: config.Get().GetApiServer().GetCors().GetAllowCredentials(),
	}).Handler(h2c.NewHandler(mux, &http2.Server{}))
	svr := http.Server{
		Addr:         fmt.Sprintf(":%d", config.Get().GetApiServer().GetPort()),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	errChan := make(chan error)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		logger.Notice(ctx, "Start server.")
		if err := svr.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()
	defer func() {
		if err := svr.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case err := <-errChan:
		logger.Emergency(ctx, map[string]interface{}{
			"message": "Failed to start server.",
			"error":   err.Error(),
		})
		cancel()
	case <-sigChan:
		logger.Notice(ctx, "Shutdown signal received, shutting down process...")
	}
}
