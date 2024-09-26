package ctxval

import (
	"context"
)

type (
	traceIDKey struct{}
)

func GetTraceID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(traceIDKey{}).(string)
	return v, ok
}

func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}
