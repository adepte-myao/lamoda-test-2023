package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adepte-myao/lamoda-test-2023/configs"
	"github.com/adepte-myao/lamoda-test-2023/internal/pkg/logger"
)

type Server struct {
	http   *http.Server
	logger logger.Logger
}

func New(config configs.AppConfig, handler http.Handler, logger logger.Logger) *Server {
	addr := fmt.Sprintf("%s:%d", config.Server.ListenAddr, config.Server.Port)

	return &Server{
		http: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  time.Second * time.Duration(config.Server.ReadTimeoutSeconds),
			WriteTimeout: time.Second * time.Duration(config.Server.WriteTimeoutSeconds),
			IdleTimeout:  time.Second * time.Duration(config.Server.IdleTimeoutSeconds),
		},
		logger: logger,
	}
}

func (server *Server) Run() error {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error)

	go func() {
		errChan <- server.http.ListenAndServe()
	}()

	select {
	case sig := <-sigChan:
		server.logger.Info("Received terminate, graceful shutdown. Signal: ", sig)

		terminationCtx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		err := server.http.Shutdown(terminationCtx)
		if err != nil {
			err = fmt.Errorf("server termination: %w", err)
			server.logger.Error(err)
		}
	case err := <-errChan:
		return fmt.Errorf("unexpected server error: %w", err)
	}

	return nil
}
