package rwriter

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	inner http.ResponseWriter
	Body  *bytes.Buffer
}

func NewWriter(inner http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{inner, new(bytes.Buffer)}
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.Body.Write(b)
	return rw.inner.Write(b)
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.inner.Header()
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.inner.WriteHeader(statusCode)
}
