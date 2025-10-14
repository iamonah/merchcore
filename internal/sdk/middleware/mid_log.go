package middleware

import (
	"net/http"
)

//contains a wrapped reponsewriter so we can capture statuscode

type RequestID struct{}

var RequestIdKey = RequestID{}

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode    int
	HeaderWritten bool
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
	}
}

func (m *ResponseWriter) Write(b []byte) (int, error) {
	if !m.HeaderWritten {
		m.StatusCode = http.StatusOK
		m.HeaderWritten = true
	}

	return m.ResponseWriter.Write(b)
}

func (m *ResponseWriter) WriteHeader(statusCode int) {
	if m.HeaderWritten {
		return
	}

	m.StatusCode = statusCode
	m.HeaderWritten = true
	m.ResponseWriter.WriteHeader(statusCode)
}

func (m *ResponseWriter) Unwrap() http.ResponseWriter {
	return m.ResponseWriter
}

// func RequestLogger(log *zerolog.Logger) base.Middleware {
// 	return func(next base.HTTPHandlerWithErr) base.HTTPHandlerWithErr {
// 		return func(w http.ResponseWriter, r *http.Request) error {
// 			start := time.Now()

// 			reqID := r.Context().Value(RequestIdKey).(string)
// 			ip := base.GetClientIP(r)

// 			nwr := NewResponseWriter(w)
// 			defer func() {
// 				event := log.Info()
// 				if nwr.statusCode >= 500 {
// 					event = log.Error()
// 				}
// 				event.
// 					Str("request_id", reqID).
// 					Str("method", r.Method).
// 					Str("url", r.URL.Path).
// 					Str("client_ip", ip).
// 					Str("user_agent", r.UserAgent()).
// 					Int("status_code", nwr.statusCode).
// 					Dur("latency", time.Since(start)).
// 					Msg("incoming request")
// 			}()

// 			err := next(nwr, r)
// 			fmt.Println(nwr.statusCode)
// 			return err
// 		}
// 	}
// }
