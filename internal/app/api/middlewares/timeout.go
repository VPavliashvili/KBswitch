package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"kbswitch/internal/core/common"
	"kbswitch/internal/core/common/middleware/models"
	"net/http"
	"time"
)

func Timeout(t int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(t)*time.Second)

			r = r.WithContext(ctx)
			rw := &models.ResponseWriterWithTimeout{ResponseWriter: w}

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

		}
		return http.HandlerFunc(fn)
	}
}
