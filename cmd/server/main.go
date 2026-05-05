package main

import (
	"context"
	"fmt"
	grpcRouter "gopher-order-service/internal/infrastructure/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"gopher-order-service/internal/application/usecases/order"
	"gopher-order-service/internal/config"
	"gopher-order-service/internal/core/ports"
	"gopher-order-service/internal/infrastructure/database"
	"gopher-order-service/internal/infrastructure/database/repositories"
	"gopher-order-service/internal/infrastructure/http"

	httpRouter "gopher-order-service/internal/presentation/http"
	"gopher-order-service/internal/presentation/http/handlers/user"
	"gopher-order-service/pkg/logger"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Core dependencies
	container.Provide(config.LoadConfig)
	container.Provide(logger.NewLogger)

	// Infrastructure
	container.Provide(database.NewPostgresDB)
	container.Provide(repositories.NewOrderPostgresRepository)
	container.Provide(func(cfg *config.Config) ports.RestaurantServiceClient {
		return http.NewRestaurantHttpClient(cfg.App.RestaurantServiceUrl)
	})

	// Application
	container.Provide(order.NewCreateOrderUseCase)

	// Presentation
	container.Provide(user.NewOrderHandler)
	container.Provide(user.NewRouter)
	container.Provide(httpRouter.NewRouter)

	return container
}

func main() {
	container := BuildContainer()

	err := container.Invoke(func(cfg *config.Config, log *zap.Logger, router *gin.Engine) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// Start HTTP Server
		httpServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.App.HTTPPort),
			Handler: router,
		}

		go func() {
			log.Info("Starting Order HTTP Server", zap.Int("port", cfg.App.HTTPPort))
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal("HTTP Server failed to start", zap.Error(err))
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the servers
		<-ctx.Done()
		log.Info("Shutting down gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error("HTTP Server forced to shutdown", zap.Error(err))
		}

		log.Info("Server exited gracefully")
	})

	if err != nil {
		fmt.Printf("Error starting application: %v\n", err)
		os.Exit(1)
	}
}
