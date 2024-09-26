package interceptor

import (
	"context"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/app/infrastructure/trace"
)

func NewTraceInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			ctx, span := trace.StartSpan(ctx, req.Spec().Procedure)
			defer span.End()

			return next(ctx, req)
		}
	}
}
