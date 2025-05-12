package mid

import (
	"context"
	"github.com/sergdort/Social/foundation/otel"
	"github.com/sergdort/Social/foundation/web"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

// Otel starts the otel tracing and stores the trace id in the context.
func Otel(tracer trace.Tracer) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			ctx = otel.InjectTracing(ctx, tracer)

			return next(ctx, r)
		}

		return h
	}

	return m
}
