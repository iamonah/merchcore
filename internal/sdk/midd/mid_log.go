package midd

import (
	"net/http"
)

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
