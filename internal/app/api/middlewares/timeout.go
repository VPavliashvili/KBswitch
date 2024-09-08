package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"kbswitch/internal/core/common"
	"net/http"
	"time"
)

const timeout = 3 // will go to config

type responseWriterWithTimeout struct {
	http.ResponseWriter
	headerWritten bool
}

func (rw *responseWriterWithTimeout) WriteHeader(statusCode int) {
	if !rw.headerWritten {
		rw.ResponseWriter.WriteHeader(statusCode)
		rw.headerWritten = true
	}
}

func (rw *responseWriterWithTimeout) Write(b []byte) (int, error) {
	if !rw.headerWritten {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func Timeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout*time.Second))

		r = r.WithContext(ctx)
		rw := &responseWriterWithTimeout{ResponseWriter: w}

		done := make(chan struct{})
		go func() {
			next.ServeHTTP(rw, r)
			close(done)
		}()

		// this is identical to defer cancel()
		select {
		case <-ctx.Done():
			cancel()

			// timeout occured
			if ctx.Err() == context.DeadlineExceeded {
				rw.WriteHeader(http.StatusGatewayTimeout)
				e := common.APIError{
					Status:  http.StatusGatewayTimeout,
					Message: ctx.Err().Error(),
				}
				j, _ := json.Marshal(e)
				fmt.Fprint(rw, string(j[:]))
			} else {
				// have no idea how to handle otherwise
				panic("request context returned unexpected type of error\n" + ctx.Err().Error())
			}
		case <-done:
			// finish without timeout
			cancel()
		}

	})
}
