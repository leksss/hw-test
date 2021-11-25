package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/logger"
)

type ServerConf struct {
	Host string
	Port string
}

type Server struct {
	srv    *http.Server
	logger logger.Log
}

func NewServer(logger logger.Log, cfg ServerConf) *Server {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: createRouter(logger),
	}
	return &Server{
		srv:    server,
		logger: logger,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) GetServerAddr() string {
	return s.srv.Addr
}

func createRouter(logger logger.Log) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "<h1>Hello</h1><div>Index page for Calendar</div>")
	})
	return loggingMiddleware(mux, logger)
}
