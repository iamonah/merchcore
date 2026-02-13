package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Primary       Primary             `mapstructure:"PRIMARY"`
	Server        ServerConfig        `mapstructure:"SERVER"`
	Database      DatabaseConfig      `mapstructure:"DATABASE"`
	Auth          AuthConfig          `mapstructure:"AUTH"`
	Redis         RedisConfig         `mapstructure:"REDIS"`
	Mailer        MailerConfig        `mapstructure:"MAILER"`
	Observability ObservabilityConfig `mapstructure:"OBSERVABILITY"`
	AWSS3         AWSS3Config         `mapstructure:"AWSS3"`
	HealthChecks  HealthChecksConfig  `mapstructure:"HEALTH_CHECKS"`
	Logging       LoggingConfig       `mapstructure:"LOGGING"`
	NewRelic      NewRelicConfig      `mapstructure:"NEW_RELIC"`
}

type Primary struct {
	Env string `mapstructure:"ENV" validate:"required"`
}

type ServerConfig struct {
	ReadTimeout        time.Duration `mapstructure:"READ_TIMEOUT" validate:"required"`
	WriteTimeout       time.Duration `mapstructure:"WRITE_TIMEOUT" vdalidate:"required"`
	IdleTimeout        time.Duration `mapstructure:"IDLE_TIMEOUT" validate:"required"`
	Port               string        `mapstructure:"PORT" validate:"required"`
	CORSAllowedOrigins []string      `mapstructure:"CORS_ALLOWED_ORIGINS" validate:"required"`
}

type DatabaseConfig struct {
	ConnMaxLifetime time.Duration `mapstructure:"CONN_MAX_LIFETIME" validate:"required"`
	ConnMaxIdleTime time.Duration `mapstructure:"CONN_MAX_IDLE_TIME" validate:"required"`
	MaxConns        int           `mapstructure:"MAX_CONNS" validate:"required"`
	MinConns        int           `mapstructure:"MIN_CONNS" validate:"required"`
	Port            int           `mapstructure:"PORT" validate:"required"`
	Host            string        `mapstructure:"HOST" validate:"required"`
	User            string        `mapstructure:"USER" validate:"required"`
	Password        string        `mapstructure:"PASSWORD" validate:"required"`
	Name            string        `mapstructure:"NAME" validate:"required"`
	SSLMode         string        `mapstructure:"SSL_MODE" validate:"required"`
}

type RedisConfig struct {
	Address  string `mapstructure:"ADDRESS" validate:"required"`
	Password string `mapstructure:"PASSWORD" validate:"required"`
}

type MailerConfig struct {
	APIKey       string `mapstructure:"MAILER_API_KEY" validate:"required"`
	Sender       string `mapstructure:"MAILER_SENDER" validate:"required,email"`
	SMTPHost     string `mapstructure:"MAILER_SMTP_HOST" validate:"required"`
	SMTPUser     string `mapstructure:"MAILER_SMTP_USER" validate:"required"`
	SMTPPassword string `mapstructure:"MAILER_SMTP_PASSWORD" validate:"required"`
	SMTPPort     int    `mapstructure:"MAILER_SMTP_PORT" validate:"required"`
}

type AuthConfig struct {
	AccessTokenLifeTime  time.Duration `mapstructure:"ACCESS_TOKEN_LIFETIME" validate:"required"`
	RefreshTokenLifeTime time.Duration `mapstructure:"REFRESH_TOKEN_LIFETIME" validate:"required"`
	TokenSymmetricKey    string        `mapstructure:"SYMMETRIC_KEY" validate:"required,min=24"`
	GoogleClientID       string        `mapstructure:"GOOGLE_CLIENT_ID" validate:"required"`
}

type AWSS3Config struct {
	Region          string `mapstructure:"REGION" validate:"required"`
	AccessKeyID     string `mapstructure:"ACCESS_KEY_ID" validate:"required"`
	SecretAccessKey string `mapstructure:"SECRET_ACCESS_KEY" validate:"required"`
	UploadBucket    string `mapstructure:"UPLOAD_BUCKET" validate:"required"`
	EndpointURL     string `mapstructure:"ENDPOINT_URL" validate:"required"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	// Validate special sections
	if err := cfg.Logging.Validate(); err != nil {
		return nil, fmt.Errorf("logging validation: %w", err)
	}

	if err := cfg.Observability.Validate(); err != nil {
		return nil, fmt.Errorf("observability validation: %w", err)
	}

	validate := validator.New()
	err := validate.Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
