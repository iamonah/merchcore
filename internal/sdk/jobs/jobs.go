package jobs

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/iamonah/merchcore/internal/config"
	"github.com/iamonah/merchcore/internal/domain/users"

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
	WelcomeEmailJob(user *users.User, code string) error
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
	UserID         uuid.UUID
	Email          string
	FirstName      string
	Code           string
	UnsubscribeURL string
}

// TODO: Unsubscribe feature
func (jq *JobClient) WelcomeEmailJob(user *users.User, code string) error {
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	unsubscribeToken := base64.URLEncoding.EncodeToString(tokenBytes)
	unsubscribeURL := fmt.Sprintf("https://yourapp.com/unsubscribe?token=%s&email=%s", unsubscribeToken, user.Email.Address)

	var buf bytes.Buffer
	payload := VerifyEmailPayload{
		UserID:         user.UserID,
		FirstName:      user.FirstName,
		Code:           code,
		Email:          user.Email.Address,
		UnsubscribeURL: unsubscribeURL,
	}

	if err := gob.NewEncoder(&buf).Encode(payload); err != nil {
		return fmt.Errorf("encode payload: type:%v: %w", TypeEmailVerify, err)
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(QueueCritical),
	}

	task := asynq.NewTask(TypeEmailVerify, buf.Bytes(), opts...)

	info, err := jq.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("enqueue: type:%v: %w", TypeEmailVerify, err)
	}

	jq.logger.Info().Str("task", TypeEmailVerify).Str("queue", info.Queue).
		Str("user_id", payload.UserID.String()).Msg("email verification enqueued")

	return nil
}

// for resending a an activation token
func (jq *JobClient) ResendVerificationTokenJob(firstname, code string, userID uuid.UUID) error {
	var buf bytes.Buffer
	payload := VerifyEmailPayload{
		FirstName: firstname,
		Code:      code,
	}

	if err := gob.NewEncoder(&buf).Encode(payload); err != nil {
		return fmt.Errorf("gob encode: type:%v :%w", TypeEmailVerify, err)
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(QueueCritical),
	}

	task := asynq.NewTask(TypeEmailVerify, buf.Bytes(), opts...)

	info, err := jq.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("enqueue email: type:%v, :%w", TypeEmailVerify, err)
	}

	jq.logger.Info().Str("user_id", userID.String()).Str("user_name", firstname).
		Str("task_type", TypeEmailVerify).Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).Msg("email verification enqueued")
	return nil
}

// for reseting a password
func (jq *JobClient) PasswordResetEmailJob(email string, token string, userId uuid.UUID) error {
	return nil
}
func (jq *JobClient) StoreCreationJob() error {
	return nil
}
