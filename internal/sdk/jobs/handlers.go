package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

func (rt *JobProcessor) JobSendVerifyEmail(ctx context.Context, t *asynq.Task) error {
	var payload VerifyEmailPayload
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		rt.logger.Error().
			Err(err).
			Str("user", payload.FirstName).
			Msg("JobSendVerifyEmail: failed to unmarshal payload")
		return fmt.Errorf("bad payload: %w", asynq.SkipRetry)
	}

	//Todo: send email to user
	rt.logger.Info().
		Str("type", t.Type()).
		// Str("to", user.Email.String()).
		Msg("JobSendVerifyEmail: successfully sent verification email")
	return nil
}
