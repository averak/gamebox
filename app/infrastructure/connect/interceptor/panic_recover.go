package interceptor

import (
	"context"
	"runtime"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/app/core/logger"
	"github.com/averak/gamebox/app/core/numunit"
	"github.com/averak/gamebox/app/infrastructure/connect/error_response"
	"github.com/averak/gamebox/protobuf/api/api_errors"
)

func NewPanicRecoverInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (_ connect.AnyResponse, retErr error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			panicked := true
			defer func() {
				if panicked {
					if err := recover(); err != nil {
						buf := make([]byte, numunit.KiB)
						for {
							n := runtime.Stack(buf, false)
							if n < len(buf) {
								buf = buf[:n]
								break
							}
							buf = make([]byte, 2*len(buf))
						}

						logger.Critical(ctx, map[string]any{
							"message": "panic occurred",
							"error":   err,
							"stack":   string(buf),
						})
						retErr = error_response.New(api_errors.ErrorCode_COMMON_UNKNOWN, api_errors.ErrorSeverity_ERROR_SEVERITY_ERROR, "unknown error occurred")
					}
				}
			}()
			res, err := next(ctx, req)
			panicked = false
			return res, err
		}
	}
}
