package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/config"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/logger"
	memory "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/storage/memory"
	mysql "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/storage/sql"
	internalhttp "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"go.uber.org/zap"
)

func main() {
	configFile := flag.String("config", "configs/config.yaml", "path to conf file")
	conf := config.NewConfig(*configFile)
	err := conf.Parse()
	if err != nil {
		log.Fatal(err.Error()) //nolintlint
	}

	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
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
		dsn := fmt.Sprintf("%s:%s@(%s:3306)/%s", conf.Database.User, conf.Database.Password, conf.Database.Host, conf.Database.Name)
		db, err := sqlx.ConnectContext(ctx, "mysql", dsn)
		if err != nil {
			logg.Error(fmt.Sprintf("Connect to storage failed: %s", err.Error()))
		}
		storage = mysql.New(conf.Database, db)
		defer db.Close()
	}

	_ = storage

	server := internalhttp.NewServer(logg, conf.Server)

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
	}
}
