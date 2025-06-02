package srv

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
)

type Config struct {
	Port    int
	Handler http.Handler
	Timeout time.Duration
	Writer  io.Writer
}

type Service struct {
	writer  io.Writer
	port    int
	server  *http.Server
	Handler http.Handler
}

func New(cfg Config) *Service {

	port := 80
	var handler http.Handler = mux.NewRouter()

	if cfg.Port != 0 {
		port = cfg.Port
	}

	if cfg.Handler != nil {
		handler = cfg.Handler
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		IdleTimeout:  cfg.Timeout,
	}

	return &Service{
		server:  srv,
		Handler: handler,
		writer:  cfg.Writer,
		port:    port,
	}
}

func (s *Service) Start() error {
	if s.writer != nil {
		s.writer.Write([]byte(fmt.Sprintf("Service is starting on port %d\n", s.port)))
	}

	return s.server.ListenAndServe()
}
