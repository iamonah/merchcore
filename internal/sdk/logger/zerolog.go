package logger

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"time"

	"github.com/IamOnah/storefronthq/internal/config"

	"github.com/rs/zerolog"
)

type enviroment string

const Production enviroment = "production"
const Developement enviroment = "developement"

func SetupLog(cfg *config.ObservabilityConfig, logService string) (*zerolog.Logger, error) {
	zerolog.TimeFieldFormat = time.RFC3339

	var env enviroment = enviroment(cfg.Environment)
	level := initLevel(cfg.GetLogLevel())
	zerolog.SetGlobalLevel(level)

	var output io.Writer
	if level.String() == "info" && env == Production {
		//logService
	} else {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}

	github.com / IamOnah / storefronthqInfo, ok := debug.Readgithub.com / IamOnah / storefronthqInfo()
	if ok {
		return nil, fmt.Errorf("no github.com/IamOnah/storefronthqinfo")
	}

	logger := zerolog.New(output).With().Timestamp().Logger()
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Int("pid", os.Getpid()).
			Str("go_version", github.com/IamOnah/storefronthqInfo.GoVersion).
			Str("service", cfg.ServiceName).
			Str("environment", string(env))
	})

	// Include stack traces for errors in development
	if !cfg.IsProduction() {
		logger = logger.With().Stack().Logger()
	}

	return &logger, nil
}

func initLevel(level string) zerolog.Level {
	var logLevel zerolog.Level

	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	case "fatal":
		logLevel = zerolog.ErrorLevel
	}

	return logLevel
}
