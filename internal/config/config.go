package config

import (
	"fmt"
	"strings"
	"time"

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
	WriteTimeout       time.Duration `mapstructure:"WRITE_TIMEOUT" validate:"required"`
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
	ResendAPIKey string `mapstructure:"RESEND_API_KEY" validate:"required"`
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

// transformFlatToNested will scan viper.AllKeys() and for any key that contains
// an underscore it will split at the first underscore and re-set a nested key in viper.
// Example:
//
//	observability_new_relic_enabled -> viper.Set("observability.new_relic_enabled", value)
func transformFlatToNested() {
	keys := viper.AllKeys()

	for _, k := range keys {
		if strings.Contains(k, ".") || !strings.Contains(k, "_") {
			continue
		}

		parts := strings.SplitN(k, "_", 2)
		if len(parts) != 2 {
			continue
		}

		section := parts[0]
		rest := parts[1]

		nestedKey := fmt.Sprintf("%s.%s", section, rest)

		val := viper.Get(k)
		viper.Set(nestedKey, val)
	}
}

func LoadConfig(path string) (*Config, error) {
	var config Config

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("reading config file: %v", err)
	}

	transformFlatToNested()

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	// for _, key := range viper.AllKeys() {
	// 	fmt.Printf("%s = %v\n", key, viper.Get(key))
	// }

	// fmt.Printf("%+v", config)
	// validate := validator.New()

	// err = validate.Struct(config)
	// if err != nil {
	// 	return nil, fmt.Errorf("config validation failed: %w", err)
	// }

	// if config.Observability == nil {
	// 	config.Observability = DefaultObservabilityConfig()
	// }
	return &config, nil
}
