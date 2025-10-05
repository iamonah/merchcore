package config

import (
	"fmt"
	"time"
)

type ObservabilityConfig struct {
	HealthChecks HealthChecksConfig `mapstructure:"health_checks" validate:"required"`
	Logging      LoggingConfig      `mapstructure:"logging" validate:"required"`
	NewRelic     NewRelicConfig     `mapstructure:"new_relic" validate:"required"`
	ServiceName  string             `mapstructure:"service_name" validate:"required"`
	Environment  string             `mapstructure:"environment" validate:"required"`
}

type HealthChecksConfig struct {
	Checks   []string      `mapstructure:"checks"`
	Interval time.Duration `mapstructure:"interval" validate:"min=1s"`
	Timeout  time.Duration `mapstructure:"timeout" validate:"min=1s"`
	Enabled  bool          `mapstructure:"enabled"`
}

type LoggingConfig struct {
	Level              string        `mapstructure:"level" validate:"required"`
	SlowQueryThreshold time.Duration `mapstructure:"slow_query_threshold"`
}

type NewRelicConfig struct {
	LicenseKey                string `mapstructure:"license_key" validate:"required"`
	AppLogForwardingEnabled   bool   `mapstructure:"app_log_forwarding_enabled"`
	DistributedTracingEnabled bool   `mapstructure:"distributed_tracing_enabled"`
	DebugLogging              bool   `mapstructure:"debug_logging"`
}

func DefaultObservabilityConfig() *ObservabilityConfig {
	return &ObservabilityConfig{
		ServiceName: "",
		Environment: "development",
		Logging: LoggingConfig{
			Level:              "info",
			SlowQueryThreshold: 100 * time.Millisecond,
		},
		NewRelic: NewRelicConfig{
			LicenseKey:                "",
			AppLogForwardingEnabled:   true,
			DistributedTracingEnabled: true,
			DebugLogging:              false, // Disabled by default to avoid mixed log formats
		},
		HealthChecks: HealthChecksConfig{
			Enabled:  true,
			Interval: 30 * time.Second,
			Timeout:  5 * time.Second,
			Checks:   []string{"database", "redis"},
		},
	}
}

func (c *ObservabilityConfig) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}

	// Validate log level
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("invalid logging level: %s (must be one of: debug, info, warn, error)", c.Logging.Level)
	}

	// Validate slow query threshold
	if c.Logging.SlowQueryThreshold < 0 {
		return fmt.Errorf("logging slow_query_threshold must be non-negative")
	}

	return nil
}

func (c *ObservabilityConfig) GetLogLevel() string {
	return c.Logging.Level
}

func (c *ObservabilityConfig) IsProduction() bool {
	return c.Environment == "production"
}
