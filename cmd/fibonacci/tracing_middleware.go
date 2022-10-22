package main

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func initTracing(logger *zap.Logger) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := cfg.InitGlobalTracer(*serviceName)
	if err != nil {
		logger.Fatal("Cannot init tracing", zap.Error(err))
	}
}

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span, ctx := opentracing.StartSpanFromContext(ctx, "http request "+r.URL.Path)
		defer span.Finish()

		wrapper := NewResponseWrapper(w)

		r = r.WithContext(ctx)

		if spanContext, ok := span.Context().(jaeger.SpanContext); ok {
			w.Header().Add("x-trace-id", spanContext.TraceID().String())
		}

		next.ServeHTTP(wrapper, r)
	})
}
