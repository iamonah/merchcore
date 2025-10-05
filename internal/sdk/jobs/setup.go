package jobs

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type JobServicer interface {
	Start() error
	JobSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type JobProcessor struct {
	server *asynq.Server
	logger *zerolog.Logger
}

func NewJobProcessor(redisOpt asynq.RedisClientOpt, logger *zerolog.Logger) *JobProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
		},
	)
	return &JobProcessor{server: server, logger: logger}
}

func (js *JobProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailVerify, js.JobSendVerifyEmail)

	return js.server.Run(mux)
}

func RunJobService(redisOpts asynq.RedisClientOpt, logger *zerolog.Logger) {
	jobProcessor := NewJobProcessor(redisOpts, logger)
	defer jobProcessor.server.Stop() // ensure workers stop on shutdown

	go func() {
		jobProcessor.logger.Info().Msg("starting task processor")
		err := jobProcessor.Start()
		if err != nil {
			jobProcessor.logger.Fatal().Err(err).Msg("cannot create redis server")
		}
	}()
}
