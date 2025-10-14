package main

import (
	"github.com/IamOnah/storefronthq/internal/app/auth"
	"github.com/IamOnah/storefronthq/internal/config"
	"github.com/IamOnah/storefronthq/internal/domain/users"
	"github.com/IamOnah/storefronthq/internal/domain/users/userdb"
	"github.com/IamOnah/storefronthq/internal/infra/database"
	"github.com/IamOnah/storefronthq/internal/sdk/authz"
	"github.com/IamOnah/storefronthq/internal/sdk/jobs"
	"github.com/IamOnah/storefronthq/internal/sdk/logger"
	"github.com/IamOnah/storefronthq/internal/sdk/mailer"
	transport "github.com/IamOnah/storefronthq/internal/transport/http"
	"github.com/IamOnah/storefronthq/internal/transport/http/router"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configurations")
	}
	log.Info().Msg("configurations succesfully initialized")

	logger, err := logger.SetupLog(cfg, cfg.Observability.ServiceName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup logging:")
	}
	log.Info().Msg("logging system initialized successfully")

	dbClient, err := database.NewDB(cfg, logger)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}
	logger.Info().Msg("database connection successfull")

	defer dbClient.Close()
	err = dbClient.MigrateUp()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}
	logger.Info().Msg("database migrations successfull")

	jwtMaker := authz.NewJWTMaker(cfg.Auth.TokenSymmetricKey)
	redisClient := jobs.NewJobClient(cfg.Redis, logger)
	mailer := mailer.NewMailTrap("onahvictorc@gmail.com", "sandbox.smtp.mailtrap.io", "66333baed34748", "9f0069b4e70329", 2525)

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
		log.Fatal().Err(err).Msg("failed to initialize userService")
	}

	mux := router.SetupRouter(userService, logger, &jwtMaker)

	go func() {
		if err := jobs.RunJobService(cfg.Redis, logger, mailer); err != nil {
			logger.Fatal().Err(err).Msg("cannot create redis server")
		}
	}()

	logger.Info().Str("port", cfg.Server.Port).Msg("server started")
	err = transport.StartServer(cfg, mux, logger)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

}
