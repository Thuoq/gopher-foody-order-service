package main

import (
	"context"
	"fmt"
	"gopher-order-service/internal/core/ports"
	externalHttp "gopher-order-service/internal/infrastructure/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"go.uber.org/zap"

	"gopher-order-service/internal/application/usecases/order"
	"gopher-order-service/internal/config"

	"gopher-order-service/internal/infrastructure/database"
	"gopher-order-service/internal/infrastructure/database/repositories"
	internalHttp "gopher-order-service/internal/presentation/http"
	userInternalHttp "gopher-order-service/internal/presentation/http/handlers/user"
	"gopher-order-service/pkg/logger"
	"gopher-order-service/internal/application/saga"
	"gopher-order-service/internal/infrastructure/messaging"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Core dependencies
	container.Provide(config.LoadConfig)
	container.Provide(logger.NewLogger)

	// Infrastructure
	container.Provide(database.NewPostgresDB)
	container.Provide(repositories.NewOrderPostgresRepository)
	container.Provide(repositories.NewOutboxPostgresRepository)
	container.Provide(func(config *config.Config) (ports.RestaurantServiceClient, error) {
		return externalHttp.NewRestaurantHttpClient(config.App.RestaurantServiceUrl), nil
	})
	container.Provide(func(config *config.Config, log *zap.Logger) *messaging.KafkaPublisher {
		return messaging.NewKafkaPublisher(config.App.KafkaBrokers, log)
	})

	// The Orchestrator will use the OutboxPublisher to be atomic
	container.Provide(func(repo ports.IOutboxRepository) ports.IMessagePublisher {
		return messaging.NewOutboxPublisher(repo)
	})

	container.Provide(saga.NewOrderCreationSaga)

	container.Provide(func(cfg *config.Config, orch ports.ISagaOrchestrator, log *zap.Logger) *messaging.KafkaConsumer {
		return messaging.NewKafkaConsumer(cfg.App.KafkaBrokers, orch, log)
	})

	container.Provide(func(repo ports.IOutboxRepository, kafkaPub *messaging.KafkaPublisher, log *zap.Logger) *messaging.OutboxProcessor {
		return messaging.NewOutboxProcessor(repo, kafkaPub, log)
	})

	// Application
	container.Provide(order.NewCreateOrderUseCase)
	// Handler
	container.Provide(userInternalHttp.NewOrderHandler)
	// Presentation
	container.Provide(userInternalHttp.NewRouter)
	container.Provide(internalHttp.NewRouter)

	return container
}

func main() {
	container := BuildContainer()

	err := container.Invoke(func(
		cfg *config.Config,
		log *zap.Logger,
		router *gin.Engine,
		sagaWorker *messaging.KafkaConsumer,
		outboxWorker *messaging.OutboxProcessor,
	) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// Start Saga Background Worker
		go sagaWorker.Start(ctx)

		// Start Outbox Processor Worker
		go outboxWorker.Start(ctx)

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
