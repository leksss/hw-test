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

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/logger"
	memory "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/storage/memory"
	mysql "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/storage/sql"
	internalhttp "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/server/http"
)

func main() {
	configFile := flag.String("config", "configs/config.yaml", "path to config file")

	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := NewConfig(*configFile)
	err := config.Parse()
	if err != nil {
		log.Fatal(err.Error()) //nolint
	}

	logg := logger.New(config.Logger, config.projectRoot)

	storage := createStorageInstance(ctx, config, logg)
	defer storage.Close(ctx)

	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar, config.Server)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info(fmt.Sprintf("calendar is running on %s", server.GetServerAddr()))
	if err := server.Start(); errors.Is(err, http.ErrServerClosed) {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolintlint
	}
}

func createStorageInstance(ctx context.Context, config Config, logg logger.Log) interfaces.Storage {
	var storage interfaces.Storage
	if config.Env == EnvTest {
		storage = memory.New()
	} else {
		storage = mysql.New(config.Database)
	}

	if err := storage.Connect(ctx); err != nil {
		logg.Error(fmt.Sprintf("Connect to storage failed: %s", err.Error()))
	}
	return storage
}
