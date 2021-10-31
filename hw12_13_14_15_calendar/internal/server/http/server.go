package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/app"
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

func NewServer(logger logger.Log, app *app.App, cfg ServerConf) *Server {
	router := createRouter()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: logRequest(router, logger),
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

func logRequest(handler http.Handler, logger logger.Log) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := &responseObserver{ResponseWriter: w}
		handler.ServeHTTP(o, r)
		addr := r.RemoteAddr
		if i := strings.LastIndex(addr, ":"); i != -1 {
			addr = addr[:i]
		}
		logger.Info(fmt.Sprintf("%s - - [%s] %q %d %d %q %q",
			addr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto),
			o.status,
			o.written,
			r.Referer(),
			r.UserAgent()))
	})
}

func createRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "<h1>Hello</h1><div>Index page for Calendar</div>")
	})
	return mux
}

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}
