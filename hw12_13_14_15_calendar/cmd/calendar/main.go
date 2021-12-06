package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/config"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/logger"
	memory "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/storage/memory"
	mysql "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/storage/sql"
	grpc "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/server/grpc"
	http "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate ./proto_generator.sh

const (
	appShutdownMessage      = "calendar exits"
	gracefulShutdownTimeout = 3 * time.Second
)

func main() {
	configFile := flag.String("config", "configs/config.yaml", "path to conf file")
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf := config.NewConfig(*configFile)
	err := conf.Parse()
	if err != nil {
		log.Fatal(err.Error()) //nolintlint
	}

	var zapConfig zap.Config
	if conf.IsDebug() {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	logg := logger.New(zapConfig, conf.Logger, conf.GetProjectRoot())

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage interfaces.Storage
	if conf.Env == config.EnvTest {
		storage = memory.New()
	} else {
		db, err := sqlx.ConnectContext(ctx, "mysql", conf.Database.DSN())
		if err != nil {
			logg.Error(fmt.Sprintf("Connect to storage failed: %s", err.Error()))
		}
		storage = mysql.New(db, logg)
		defer db.Close()
	}

	grpcServer := grpc.NewServer(logg, conf, storage)
	httpServer := http.NewServer(logg, conf, storage)
	errs := make(chan error)

	go func() {
		errs <- grpcServer.StartGRPC()
	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
		errs <- httpServer.StartHTTPProxy()
	}()

	go func() {
		<-ctx.Done()

		success := make(chan string)
		go func() {
			errs <- httpServer.StopHTTPProxy(context.Background())
			success <- "HTTP server successfully stopped"
		}()
		go func() {
			grpcServer.StopGRPC()
			success <- "gRPC server successfully stopped"
		}()
		go func() {
			time.Sleep(gracefulShutdownTimeout)
			logg.Error("failed to gracefully shut down server within timeout. Shutting down with Fatal",
				zap.Duration("timeout", gracefulShutdownTimeout))
		}()
		logg.Info(<-success)
		logg.Info(<-success)
		errs <- errors.New(appShutdownMessage)
	}()

	for err := range errs {
		if err == nil {
			continue
		}
		logg.Warn("shutdown err message", zap.Error(err))
		if err.Error() == appShutdownMessage {
			return
		}
	}
}
