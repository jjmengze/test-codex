package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log-receiver/internal/handler"
	"log-receiver/internal/repo"
	"log-receiver/internal/usecase"
	"log-receiver/pkg/aws"
	"log-receiver/pkg/aws/kinesis"
	logger "log-receiver/pkg/logger"
	slogger "log-receiver/pkg/logger/slog"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	gin.SetMode(gin.DebugMode)
	logger := slogger.GetGlobalLogger()

	//inject config(e.g. aws service )
	kinesisClient := kinesis.NewClient(aws.NewKinesisClient())

	engine := gin.New()
	engine.Use(gin.Recovery())

	//inject repo
	repoPublisher := repo.NewPublisher(logger, kinesisClient)
	//inject usecase

	usecaseReceiver := usecase.NewReceiver(logger, repoPublisher)
	usecaseValidator := usecase.NewValidator(logger)
	//inject handler
	_ = handler.NewHttpHandler(logger, engine, usecaseReceiver, usecaseValidator)

	startWithGracefulShutdown(logger, engine, 1*time.Minute)

}

func startWithGracefulShutdown(log logger.Logger, router *gin.Engine, duration time.Duration) {
	h2s := &http2.Server{
		MaxConcurrentStreams: 200,
		IdleTimeout:          60 * time.Second,
	}
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      h2c.NewHandler(router, h2s),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 150 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.FatalF("Start Server Failed: %v", err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.InfoF("Shutdown Server ... wait for %v seconds", duration.Seconds())

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.FatalF("Server Shutdown Failed: %v", err)
	}

	log.InfoW("Server Exiting")
}
