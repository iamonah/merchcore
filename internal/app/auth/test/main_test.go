package auth_test

import (
	"os"
	"testing"

	"github.com/iamonah/merchcore/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

var (
	cfg  *config.Config
	log  zerolog.Logger
	pool *pgxpool.Pool
)

func TestMain(m *testing.M) {
	cfg = &config.Config{}
	var err error
	for _, p := range []string{".", "../../../.."} {
		cfg, err = config.LoadConfig(p)
		if err == nil {
			break
		}
	}

	log = zerolog.New(os.Stdout)

	os.Exit(m.Run())
}
