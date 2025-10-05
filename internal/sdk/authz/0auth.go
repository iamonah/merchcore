package authz

import (
	"context"

	"google.golang.org/api/idtoken"
)

func VerifyGoogleIDToken(ctx context.Context, token string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(ctx, token, "YOUR_GOOGLE_CLIENT_ID")
	if err != nil {
		return nil, err
	}
	return payload, nil
}
