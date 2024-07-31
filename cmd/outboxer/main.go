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

	"github.com/gin-gonic/gin"
	"github.com/mirhijinam/outboxer/internal/api"
	"github.com/mirhijinam/outboxer/internal/config"
	"github.com/mirhijinam/outboxer/internal/pkg/db"
	"github.com/mirhijinam/outboxer/internal/pkg/logger"
	"github.com/mirhijinam/outboxer/internal/repository"
	"github.com/mirhijinam/outboxer/internal/service/eventhandler"
	"github.com/mirhijinam/outboxer/internal/service/kafka"
	"github.com/mirhijinam/outboxer/internal/service/message"
	"go.uber.org/zap"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal(err.Error())
	}

	pool, err := db.MustOpenDB(context.Background(), config.DBConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	lgr := logger.New(config.LoggerConfig.Mode, config.LoggerConfig.Filepath)
	defer lgr.Sync()

	lgr.Info("this is an info message")

	rep := repository.New(pool)
	r, err := api.New(
		message.New(rep),
		lgr,
	)
	if err != nil {
		log.Fatalln(err)
	}

	kafkaProducer, err := kafka.NewProducer(config.KafkaConfig, lgr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer kafkaProducer.Close()

	eventHandler := eventhandler.New(config.EventHandlerConfig, rep, kafkaProducer, lgr)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventHandler.StartHandlingEvents(ctx)

	kafkaConsumer, err := kafka.NewConsumer(config.KafkaConfig, lgr, func(ctx context.Context, message []byte) error {
		fmt.Printf("received message: %s\n", string(message))
		return nil
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer kafkaConsumer.Close()

	go kafkaConsumer.Run(ctx)

	err = run(r, config.ServerConfig)
	if err != nil {
		lgr.Error("failed to run the server",
			zap.String("error", err.Error()),
		)
		log.Println(err)
		os.Exit(2)
	}
}

func run(r *gin.Engine, srvCfg config.ServerConfig) error {
	timeout, err := time.ParseDuration(srvCfg.Timeout)
	if err != nil {
		return fmt.Errorf("failed to parse timeout value: %w", err)
	}
	idletimeout, err := time.ParseDuration(srvCfg.Timeout)
	if err != nil {
		return fmt.Errorf("failed to parse idle timeout value: %w", err)
	}

	srv := http.Server{
		Handler:      r,
		Addr:         ":" + srvCfg.Port,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
		IdleTimeout:  idletimeout,
	}

	serveChan := make(chan error, 1)
	go func() {
		serveChan <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	case err := <-serveChan:
		return err
	}
}
