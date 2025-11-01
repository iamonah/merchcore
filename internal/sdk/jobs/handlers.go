package jobs

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	UserWelcomeTemplate = "welcomemail.html"
)

func (rt *JobProcessor) DoWelcomeEmailJob(ctx context.Context, t *asynq.Task) error {
	var payload VerifyEmailPayload
	if err := gob.NewDecoder(bytes.NewReader(t.Payload())).Decode(&payload); err != nil {
		rt.logger.Error().Err(err).
			Str("user", payload.FirstName).
			Str("type", t.Type()).
			Str("user_id", payload.UserID.String()).
			Msg("decode failed")
		return fmt.Errorf("gob decode: %w: %w", asynq.SkipRetry, err)
	}

	retryCount, _ := asynq.GetRetryCount(ctx)

	if err := rt.mailer.Send(UserWelcomeTemplate, payload.Email, payload); err != nil {
		rt.logger.Error().Err(err).
			Str("type", t.Type()).
			Str("email", payload.Email).
			Int("attempt", retryCount).
			Msg("send failed")
		return fmt.Errorf("send email: %w", err)
	}

	rt.logger.Info().Str("type", t.Type()).Str("email", payload.Email).
		Int("attempt", retryCount).Msg("email sent")
	return nil
}
