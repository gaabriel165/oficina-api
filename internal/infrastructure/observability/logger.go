package observability

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type requestIDKey struct{}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey{}).(string)
	if !ok {
		return ""
	}
	return requestID
}

type correlationHandler struct {
	slog.Handler
}

func (h correlationHandler) Handle(ctx context.Context, record slog.Record) error {
	if requestID := RequestIDFromContext(ctx); requestID != "" {
		record.AddAttrs(slog.String("request_id", requestID))
	}

	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.HasTraceID() {
		record.AddAttrs(
			slog.String("trace_id", spanContext.TraceID().String()),
			slog.String("span_id", spanContext.SpanID().String()),
		)
	}

	return h.Handler.Handle(ctx, record)
}

func (h correlationHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return correlationHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h correlationHandler) WithGroup(name string) slog.Handler {
	return correlationHandler{Handler: h.Handler.WithGroup(name)}
}

func NewLogger(serviceName, environment string) *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(correlationHandler{Handler: jsonHandler}).With(
		slog.String("service.name", serviceName),
		slog.String("environment", environment),
	)
	slog.SetDefault(logger)
	return logger
}

func Event(ctx context.Context, name string, attrs ...any) {
	slog.Default().InfoContext(ctx, name, append([]any{slog.String("event", name)}, attrs...)...)
}

func Failure(ctx context.Context, name string, err error, attrs ...any) {
	fields := append([]any{slog.String("event", name), slog.String("error", err.Error())}, attrs...)
	slog.Default().ErrorContext(ctx, name, fields...)
}
