package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go-boilerplate/internal/configure"
	examplegateway "go-boilerplate/internal/gateway/example"
	httpexample "go-boilerplate/internal/handlers/http/example"
	httphealth "go-boilerplate/internal/handlers/http/health"
	"go-boilerplate/internal/metrics"
	presenterexample "go-boilerplate/internal/presenter/example"
	examplerepository "go-boilerplate/internal/repository/example"
	usecaseexample "go-boilerplate/internal/usecase/example"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := configure.MustLog()
	db := configure.MustNewDB()
	defer db.Close()

	exampleRepo := examplerepository.NewExamplePG(db)
	exampleGateway := examplegateway.NewExampleHTTP(os.Getenv("EXAMPLE_API_BASE_URL"))

	if err := exampleGateway.Health(ctx); err != nil {
		logger.WarnContext(ctx, "example gateway health check failed", "error", err.Error())
	}

	exampleUsecase := usecaseexample.New(exampleRepo)
	examplePresenter := presenterexample.NewPresenter()

	exampleHandler := httpexample.NewHandler(exampleUsecase, examplePresenter, logger)
	healthHandler := httphealth.NewHandler()
	appMetrics := metrics.New()

	router := gin.Default()
	router.Use(appMetrics.Middleware())
	router.GET("/_health", healthHandler.Handle)
	router.GET("/metrics", gin.WrapH(appMetrics.Handler()))

	api := router.Group("/api")
	api.GET("/examples", exampleHandler.List)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "9000"
	}

	logger.InfoContext(ctx, "http server started", "port", port)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.ErrorContext(ctx, "http server stopped with error", "error", err.Error())
	}
}
