package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
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

	awsSDK "github.com/aws/aws-sdk-go/aws"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var port = flag.Int("port", 8080, "port to listen on, default 8080")
var isTestPem = flag.Bool("is_test_pem", false, "is used mock pem to verify jwt token")

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}
	flag.Parse()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gin.SetMode(gin.DebugMode)
	logger := slogger.GetGlobalLogger()

	//inject config(e.g. aws service )
	awsCfg := awsSDK.NewConfig().WithRegion("us-west-2")
	awsSession, err := aws.NewSession(awsCfg)
	ssmClient := aws.NewSsmClient(awsSession)
	//TODO: This is not a best practice because it introduces a global variable and it's hard to guarantee this function is called only once, even with sync.Once.
	aws.InitPublicKeyMap(ssmClient)

	kinesisClient, err := kinesis.NewClient(ctx, logger, aws.NewKinesisClient(awsSession))
	if err != nil {
		logger.WithContext(ctx).FatalF("create kinesis client error: %v", err)
		return
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	//inject repo
	repoPublisher := repo.NewPublisher(logger, kinesisClient)
	//inject usecase
	usecaseReceiver := usecase.NewReceiver(logger, repoPublisher)
	usecaseValidator := usecase.NewValidator(logger)
	//inject handler
	_ = handler.NewHttpHandler(logger, engine, usecaseReceiver, usecaseValidator, *isTestPem)

	startWithGracefulShutdown(logger, engine, 1*time.Minute)

}

func startWithGracefulShutdown(log logger.Logger, router *gin.Engine, duration time.Duration) {
	h2s := &http2.Server{
		MaxConcurrentStreams: 200,
		IdleTimeout:          60 * time.Second,
	}
	srv := &http.Server{
		Addr:         fmt.Sprint(":", *port),
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
