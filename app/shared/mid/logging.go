package mid

import (
	"context"
	"errors"
	"fmt"
	"github.com/sergdort/Social/app/shared/errs"
	"github.com/sergdort/Social/foundation/logger"
	"github.com/sergdort/Social/foundation/web"
	"net/http"
	"time"
)

// Logger writes information about the request to the logs.
func Logger(log *logger.Logger) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			now := time.Now()

			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Info(ctx, "request started", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr)

			resp := next(ctx, r)
			err := isError(resp)

			var statusCode = errs.OK
			if err != nil {
				statusCode = errs.Internal

				var v *errs.Error
				if errors.As(err, &v) {
					statusCode = v.Code
				}
			}

			log.Info(ctx, "request completed", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr,
				"statuscode", statusCode, "since", time.Since(now).String())

			return resp
		}

		return h
	}

	return m
}

// isError tests if the Encoder has an error inside of it.
func isError(e web.Encoder) error {
	err, isError := e.(error)
	if isError {
		return err
	}
	return nil
}
