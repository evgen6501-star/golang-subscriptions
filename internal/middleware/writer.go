package middleware

import "net/http"

var (
	StatusCodeUnitialized = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUnitialized,
	}
}
func (rw *ResponseWriter) WriteHeader(statisCode int) {
	rw.ResponseWriter.WriteHeader(statisCode)
	rw.statusCode = statisCode
}
func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUnitialized {
		panic("no status code")
	}
	return rw.statusCode
}
