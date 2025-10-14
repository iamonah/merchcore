package jobs

import (
	"context"
	"fmt"

	"github.com/IamOnah/storefronthq/internal/config"
	"github.com/IamOnah/storefronthq/internal/sdk/mailer"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type JobServicer interface {
	Start() error
	DoWelcomeEmailJob(ctx context.Context, task *asynq.Task) error
}

type JobProcessor struct {
	server *asynq.Server
	logger *zerolog.Logger
	mailer *mailer.Mail
}

func NewJobProcessor(cfg config.RedisConfig, logger *zerolog.Logger, mailer *mailer.Mail) *JobProcessor {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       0,
	}
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
		},
	)
	return &JobProcessor{server: server, logger: logger, mailer: mailer}
}

func (js *JobProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailVerify, js.DoWelcomeEmailJob)

	return js.server.Run(mux)
}

func RunJobService(cfg config.RedisConfig, logger *zerolog.Logger, mailer *mailer.Mail) error {
	jobProcessor := NewJobProcessor(cfg, logger, mailer)
	defer jobProcessor.server.Stop() // ensure workers stop on shutdown

	jobProcessor.logger.Info().Msg("starting task processor")
	err := jobProcessor.Start()
	if err != nil {
		return fmt.Errorf("runjobservice: %w", err)
	}
	return nil
}
