package middlewares

import (
	"context"
	"kbswitch/internal/core/common/logger"
	"net/http"

	"github.com/google/uuid"
)

func ContentTypeJSON(next http.Handler) http.Handler {
	const (
		HeaderKeyContentType       = "Content-Type"
		HeaderValueContentTypeJSON = "application/json;charset=utf8"
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderKeyContentType, HeaderValueContentTypeJSON)
		next.ServeHTTP(w, r)
	})
}

func RequestID(next http.Handler) http.Handler {
	const (
		XRequestIDKey = "XRequestID"
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRequestID := uuid.NewString()
		w.Header().Set(XRequestIDKey, xRequestID)

		ctx := context.WithValue(r.Context(), logger.LogIDKey, xRequestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
