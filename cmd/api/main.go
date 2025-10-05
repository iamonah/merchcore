package main

import (
	"github.com/IamOnah/storefronthq/internal/app/auth"
	"github.com/IamOnah/storefronthq/internal/config"
	"github.com/IamOnah/storefronthq/internal/domain/users/userstore"
	"github.com/IamOnah/storefronthq/internal/infra/database"
	"github.com/IamOnah/storefronthq/internal/sdk/authz"
	"github.com/IamOnah/storefronthq/internal/sdk/jobs"
	"github.com/IamOnah/storefronthq/internal/sdk/logger"
	transport "github.com/IamOnah/storefronthq/internal/transport/http"
	"github.com/IamOnah/storefronthq/internal/transport/http/router"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("failed to load configurations")
	}
	log.Info().Msg("configurations succesfully initialized")

	log, err := logger.SetupLog(config.DefaultObservabilityConfig(), "storefronthq")
	if err != nil {
		log.Fatal().Msg("failed to setup logging")
	}
	log.Info().Msg("logging system initialized successfully")

	dbClient, err := database.NewDB(cfg, log)
	if err != nil {
		log.Fatal().Msg("failed to setup logging")
	}
	log.Info().Msg("database migrations successfull")

	defer dbClient.Close()
	err = dbClient.MigrateUp()
	if err != nil {
		log.Fatal().Msg("failed to setup logging")
	}
	log.Info().Msg("database migrations successfull")

	userService, err := auth.NewUserService(
		auth.WithAuth(authz.NewJWTMaker(cfg.Auth.TokenSymmetricKey)),
		auth.WithUserRepository(userstore.NewUserStore(dbClient.Pool)),
		auth.WithTrxManager(database.NewTRXManager(dbClient.Pool, log)),
		auth.WithJob(jobs.NewJobClient(cfg.Redis, log)),
		auth.WithLog(log),
	)
	if err != nil {
		log.Fatal().Msg("failed to initialize userService")
	}

	mux := router.SetupRouter(userService)

	log.Info().Str("port", cfg.Server.Port).Msg("server started")
	err = transport.StartServer(cfg, mux, log)
	if err != nil {
		log.Fatal().Msg("failed to initialize server")
	}

}
