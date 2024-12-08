package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

type Middleware struct {
	limiter *rate.Limiter
}

func NewMiddleware() *Middleware {
	return &Middleware{
		limiter: rate.NewLimiter(rate.Every(time.Second), 10), // 10 requests per second
	}
}

func (m *Middleware) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

func (m *Middleware) RequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		w.Header().Set("X-Request-ID", requestID)
		next(w, r.WithContext(ctx))
	}
}
