package middlewares

import (
	"kbswitch/internal/core/common/logger"
	"net/http"
)

func LogHttpCycle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rww := logger.NewResponseWriterWrapper(w)
		defer func() {
			msg := logger.GetRequestResponseLog(rww, r)
			logger.Trace(msg)
		}()

		next.ServeHTTP(rww, r)
	})
}
