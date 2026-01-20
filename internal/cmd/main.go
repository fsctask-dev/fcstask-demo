package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fcstask/internal/app"
	"fcstask/internal/config"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatalf("stopped: %v", err)
	}

	log.Println("app start")
}
