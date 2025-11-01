package config

import (
	"fmt"
	"time"
)

type ObservabilityConfig struct {
	ServiceName string `mapstructure:"SERVICE_NAME" validate:"required"`
	Environment string `mapstructure:"ENVIRONMENT" validate:"required"`
}

func (c *ObservabilityConfig) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}
	return nil
}

type HealthChecksConfig struct {
	Checks   []string      `mapstructure:"CHECKS"`
	Interval time.Duration `mapstructure:"INTERVAL" validate:"min=1s"`
	Timeout  time.Duration `mapstructure:"TIMEOUT" validate:"min=1s"`
	Enabled  bool          `mapstructure:"ENABLED"`
}

type LoggingConfig struct {
	Level              string        `mapstructure:"LEVEL" validate:"required"`
	SlowQueryThreshold time.Duration `mapstructure:"SLOW_QUERY_THRESHOLD"`
}

func (l *LoggingConfig) Validate() error {
	// Validate log level
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLevels[l.Level] {
		return fmt.Errorf("invalid logging level: %s (must be one of: debug, info, warn, error)", l.Level)
	}

	return nil
}

type NewRelicConfig struct {
	LicenseKey                string `mapstructure:"LICENSE_KEY" validate:"required"`
	AppLogForwardingEnabled   bool   `mapstructure:"APP_LOG_FORWARDING_ENABLED"`
	DistributedTracingEnabled bool   `mapstructure:"DISTRIBUTED_TRACING_ENABLED"`
	DebugLogging              bool   `mapstructure:"DEBUG_LOGGING"`
}

// 	// Validate slow query threshold
// 	if c.Logging.SlowQueryThreshold < 0 {
// 		return fmt.Errorf("logging slow_query_threshold must be non-negative")
// 	}

// 	return nil
// }

// func (c *ObservabilityConfig) GetLogLevel() string {
// 	return c.Logging.Level
// }

func (c *ObservabilityConfig) IsProduction() bool {
	return c.Environment == "production"
}
