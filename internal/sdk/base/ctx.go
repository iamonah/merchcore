package base

import (
	"fmt"
	"net/http"

	"github.com/iamonah/merchcore/internal/sdk/authz"
	"github.com/iamonah/merchcore/internal/sdk/midd"
)

func GetReqIDCTX(r *http.Request) (string, error) {
	v, ok := r.Context().Value(midd.RequestIdKey).(string)
	if !ok {
		return "", fmt.Errorf("reqID not in context")
	}
	return v, nil
}

func GetJWTPayloadCTX(r *http.Request) (*authz.Payload, error) {
	v, ok := r.Context().Value(midd.AuthContextPayloadKey).(*authz.Payload)
	if !ok {
		return nil, fmt.Errorf("jwt payload not in context")
	}
	return v, nil
}