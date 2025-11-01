package main

import (
	"github.com/iamonah/merchcore/internal/app/auth"
	"github.com/iamonah/merchcore/internal/config"
	"github.com/iamonah/merchcore/internal/domain/users"
	"github.com/iamonah/merchcore/internal/domain/users/userdb"
	"github.com/iamonah/merchcore/internal/infra/database"
	"github.com/iamonah/merchcore/internal/sdk/authz"
	"github.com/iamonah/merchcore/internal/sdk/jobs"
	"github.com/iamonah/merchcore/internal/sdk/logger"
	"github.com/iamonah/merchcore/internal/sdk/mailer"
	transport "github.com/iamonah/merchcore/internal/transport/http"
	"github.com/iamonah/merchcore/internal/transport/http/router"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("config load failed")
	}
	log.Info().Msg("config loaded")

	logger, err := logger.SetupLog(cfg, cfg.Observability.ServiceName)
	if err != nil {
		log.Fatal().Err(err).Msg("log setup failed")
	}
	log.Info().Msg("logger ready")

	dbClient, err := database.NewDB(cfg, logger)
	if err != nil {
		log.Fatal().Err(err).Msg("db init failed")
	}
	logger.Info().Msg("db connected")

	defer dbClient.Close()
	if err := dbClient.MigrateUp(); err != nil {
		log.Fatal().Err(err).Msg("migration failed")
	}
	logger.Info().Msg("migration done")

	jwtMaker := authz.NewJWTMaker(cfg.Auth.TokenSymmetricKey)
	redisClient := jobs.NewJobClient(cfg.Redis, logger)
	mailer := mailer.NewMailTrap(&cfg.Mailer)

	userService, err := auth.NewUserService(
		auth.WithAuth(jwtMaker),
		auth.WithUserBusiness(users.NewUserBusiness(
			userdb.Newuserdb(dbClient.Pool),
			database.NewTRXManager(dbClient.Pool, logger),
			&jwtMaker,
			cfg,
		)),
		auth.WithJob(redisClient),
		auth.WithLog(logger),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("user service init failed")
	}

	mux := router.SetupRouter(userService, logger, &jwtMaker)

	go func() {
		if err := jobs.RunJobService(cfg.Redis, logger, mailer); err != nil {
			logger.Fatal().Err(err).Msg("redis job failed")
		}
	}()

	logger.Info().Str("port", cfg.Server.Port).Msg("server starting")
	if err := transport.StartServer(cfg, mux, logger); err != nil {
		log.Fatal().Err(err).Msg("server start failed")
	}
}
