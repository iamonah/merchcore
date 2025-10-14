package jobs

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IamOnah/storefronthq/internal/config"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

const (
	TypeEmailVerify = "email:verify"
	// TypeEmailWelcome = "email:welcome"
	TypeImageResize = "image:resize"
)

type JobClient struct {
	client *asynq.Client
	logger *zerolog.Logger
}

type JobService interface {
	WelcomeEmailJob(firstname string, code string, userID uuid.UUID, email string) error
	ResendVerificationTokenJob(firstname string, code string, userID uuid.UUID) error
	PasswordResetEmailJob(email string, token string, userId uuid.UUID) error
}

func NewJobClient(cfg config.RedisConfig, logger *zerolog.Logger) *JobClient {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       0,
	}
	redisClient := asynq.NewClient(redisOpt)

	return &JobClient{client: redisClient, logger: logger}
}

func (jc *JobClient) CloseClient() {
	jc.client.Close()
}

type VerifyEmailPayload struct {
	Email          string
	FirstName      string
	Code           string
	UnsubscribeURL string
}

// TODO: remember test phase so modify the unsubscribe link
func (jq *JobClient) WelcomeEmailJob(firstname string, code string, userID uuid.UUID, email string) error {
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	unsubscribeToken := base64.URLEncoding.EncodeToString(tokenBytes)
	var unsubscribeURL = fmt.Sprintf("https://yourapp.com/unsubscribe?token=%s&email=%s", unsubscribeToken, email)

	payload, err := json.Marshal(
		VerifyEmailPayload{
			FirstName: firstname,
			Code:      code, Email: email,
			UnsubscribeURL: unsubscribeURL,
		})
	if err != nil {
		jq.logger.Error().
			Err(err).
			Str("user", firstname).
			Str("user_id", userID.String()).
			Str("task_type", TypeEmailVerify).
			Msg("failed to marshall payload")
		return fmt.Errorf("marshall payload %w", err)
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(QueueCritical),
	}

	welcomeEmailTask := asynq.NewTask(TypeEmailVerify, payload, opts...)

	info, err := jq.client.Enqueue(welcomeEmailTask)
	if err != nil {
		jq.logger.Error().
			Err(err).
			Str("queue", info.Queue).
			Str("username", firstname).
			Str("task_type", TypeEmailVerify).
			Msg("failed to enqueue email verification task")
		return fmt.Errorf("enqueue email: %w", err)
	}

	jq.logger.Info().
		Str("user", firstname).
		Str("task_type", TypeEmailVerify).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("succesfully enqueued email verification task")
	return nil
}

// for resending a an activation token
func (jq *JobClient) ResendVerificationTokenJob(firstname string, code string, userID uuid.UUID) error {
	payload, err := json.Marshal(VerifyEmailPayload{FirstName: firstname, Code: code})
	if err != nil {
		jq.logger.Error().
			Err(err).
			Str("user", firstname).
			Str("user_id", userID.String()).
			Str("task_type", TypeEmailVerify).
			Msg("failed to marshall payload")
		return fmt.Errorf("marshall payload %w", err)
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(QueueCritical),
	}

	welcomeEmailTask := asynq.NewTask(TypeEmailVerify, payload, opts...)

	info, err := jq.client.Enqueue(welcomeEmailTask)
	if err != nil {
		jq.logger.Error().
			Err(err).
			Str("queue", info.Queue).
			Str("username", firstname).
			Str("task_type", TypeEmailVerify).
			Msg("failed to enqueue email verification task")
		return fmt.Errorf("enqueue email: %w", err)
	}

	jq.logger.Info().
		Str("user", firstname).
		Str("task_type", TypeEmailVerify).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("succesfully enqueued email verification task")
	return nil
}

// for reseting a password
func (jq *JobClient) PasswordResetEmailJob(email string, token string, userId uuid.UUID) error {
	return nil
}
func (jq *JobClient) StoreCreationJob() error {
	return nil
}
