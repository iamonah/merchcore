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
	var err error
	cfg, err = config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	log = zerolog.New(os.Stdout)

	os.Exit(m.Run())
}
