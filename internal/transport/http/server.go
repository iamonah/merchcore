package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IamOnah/storefronthq/internal/config"

	"github.com/rs/zerolog"
)

func StartServer(cfg *config.Config, mux http.Handler, log *zerolog.Logger) error {
	server := http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	errHttp := make(chan error)

	go func() {
		errHttp <- server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log.Info().Str("signal", (<-quit).String()).Msg("signal recieved shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	err = <-errHttp
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server closed unexpectedly: %w", err)
	}

	log.Info().Msg("Server shutdown complete")
	return nil
}
