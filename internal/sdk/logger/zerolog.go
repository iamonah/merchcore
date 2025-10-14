package logger

import (
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

func SetupLog(cfg *config.Config, logService string) (*zerolog.Logger, error) {
	zerolog.TimeFieldFormat = time.RFC3339

	var env enviroment = enviroment(cfg.Primary.Env)
	level := initLevel(cfg.Logging.Level)
	zerolog.SetGlobalLevel(level)

	var output io.Writer
	if level.String() == "info" && env == Production {
		//logService
	} else {
		// output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output = os.Stdout
	}

	build, _ := debug.ReadBuildInfo()

	logger := zerolog.New(output).With().Timestamp().Logger()
	logger = logger.With().
		Int("pid", os.Getpid()).
		Str("go_version", build.GoVersion).
		Str("service", cfg.Observability.ServiceName).
		Str("environment", string(env)).
		Logger()

	// Include stack traces for errors in development
	if !cfg.Observability.IsProduction() {
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
