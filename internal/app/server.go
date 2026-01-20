package app

import (
	"context"
	"fmt"
	"log"

	"fcstask/internal/config"
	"fcstask/internal/server"
	"fcstask/internal/server/handler"
)

func Run(ctx context.Context, cfg *config.Config) error {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	h := handler.NewHandler()

	router := server.
		NewRouterBuilder().
		Version("/v1").
		WithEcho(h.Echo).
		Register().
		Build()

	srv := server.NewServer(addr, router)
	errCh := make(chan error, 1)

	go func() {
		log.Println("server started on", addr)
		errCh <- srv.Start(ctx)
	}()

	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		log.Println("shutdown requested")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			cfg.Server.ShutdownTimeout,
		)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return <-errCh
	}
}
