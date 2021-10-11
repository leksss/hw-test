package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
)

func main() {
	configFile := flag.String("config", "configs/config.yaml", "path to config file")

	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(*configFile)
	err := config.Parse()
	if err != nil {
		log.Fatal(err.Error())
	}

	logg := logger.New(config.Logger, config.projectRoot)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
