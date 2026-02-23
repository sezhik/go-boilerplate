package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go-boilerplate/internal/configure"
	"go-boilerplate/internal/gateway"
	"go-boilerplate/internal/repository"
	"go-boilerplate/internal/util/signal"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configure.MustLog()
	db := configure.MustNewDB()
	defer db.Close()

	repo := repository.NewExamplePG(db)
	gateway := gateway.NewExampleHTTP(os.Getenv("EXAMPLE_API_BASE_URL"))

	count, err := repo.CountExamples(ctx)
	if err != nil {
		panic("failed to count examples: " + err.Error())
	}

	slog.InfoContext(ctx, "boilerplate started", "examples_count", count)

	if err := gateway.Health(ctx); err != nil {
		slog.WarnContext(ctx, "example gateway health check failed", "error", err.Error())
	}

	signalTrap := signal.NewTrap()

	signalTrap.Wait(ctx)
	time.Sleep(time.Second)
}
