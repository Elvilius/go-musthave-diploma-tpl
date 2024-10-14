package middleware

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"net/http"
	"time"
)

type Tokener interface {
	GetUserId(tokenString string) int
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func Logging(logger zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{size: 0, status: 0}

			lw := &loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}

			h.ServeHTTP(lw, r)

			duration := time.Since(start)

			logger.Infoln(
				"uri", r.RequestURI,
				"method", r.Method,
				"status", lw.responseData.status,
				"duration", duration,
				"size", lw.responseData.size,
			)
		})
	}
}

func ValidateJWT(token Tokener) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			bearer := strings.TrimPrefix(authorization, "Bearer ")
			userID := token.GetUserId(bearer)
			if userID < 1 {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "user_id", userID))
			h.ServeHTTP(w, r)
		})
	}
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.responseData.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
