package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

func (rt *JobProcessor) DoWelcomeEmailJob(ctx context.Context, t *asynq.Task) error {
	var payload VerifyEmailPayload
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		rt.logger.Error().
			Err(err).
			Str("user", payload.FirstName).
			Msg("JobSendVerifyEmail: failed to unmarshal payload")
		return fmt.Errorf("bad payload: %w", asynq.SkipRetry)
	}

	err = rt.mailer.Send("welcomemail.html", payload.Email, payload)
	if err != nil {
		return fmt.Errorf("sendwelcomeemail: %w", err)
	}
	rt.logger.Info().
		Str("type", t.Type()).
		Str("to", payload.Email).
		Msg("dowelcomeemailjob: successfully sent verification email")
	return nil
}
