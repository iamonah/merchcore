package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/IamOnah/storefronthq/internal/sdk/base"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type RequestID struct{}

var RequestIdKey = RequestID{}

type newResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func NewResponseWriter(w http.ResponseWriter) *newResponseWriter {
	return &newResponseWriter{
		ResponseWriter: w,
	}
}

func (m *newResponseWriter) Write(b []byte) (int, error) {
	if !m.headerWritten {
		m.statusCode = http.StatusOK
		m.headerWritten = true
	}

	return m.ResponseWriter.Write(b)
}

func (m *newResponseWriter) WriteHeader(statuscode int) {
	m.ResponseWriter.WriteHeader(statuscode)

	if !m.headerWritten {
		m.statusCode = statuscode
		m.headerWritten = true
	}
}

func (m *newResponseWriter) Unwrap() http.ResponseWriter {
	return m.ResponseWriter
}

func RequestLogger(log *zerolog.Logger) base.Middlware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.New().String()
			}

			ctx := context.WithValue(r.Context(), RequestIdKey, reqID)
			nwr := NewResponseWriter(w)
			next(nwr, r.WithContext(ctx))

			ip := base.GetClientIP(r)

			event := log.Info()

			if nwr.statusCode >= 500 {
				event = log.Error()
			}

			event.
				Str("request_id", reqID).
				Str("method", r.Method).
				Str("url", r.URL.Path).
				Str("query", r.URL.RawQuery).
				Str("client_ip", ip).
				Str("user_agent", r.UserAgent()).
				Int("status_code", nwr.statusCode).
				Dur("latency", time.Since(start)).
				Msg("incoming request")
		}
	}
}
