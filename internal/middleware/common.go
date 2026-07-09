package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/evgen6501-star/golang-subscriptions/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	reqIDHeader = "X-Request-ID"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(reqIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}
			r.Header.Set(reqIDHeader, requestID)
			w.Header().Set(reqIDHeader, requestID)
			next.ServeHTTP(w, r)
		})
	}
}
func Logger(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get(reqIDHeader)
			l := log.With(
				zap.String("request_id", reqID),
				zap.String("url", r.URL.String()),
			)
			ctx := context.WithValue(r.Context(), "log", l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := NewResponseWriter(w)
			ctx := r.Context()
			log := logger.FromContext(ctx)
			before := time.Now()
			log.Debug(
				">>> incomig HTTP request",
				zap.Time("time", before.UTC()),
			)
			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status code", rw.GetStatusCodeOrPanic()),
				zap.Duration("latency", time.Now().Sub(before)),
			)
		})
	}
}
