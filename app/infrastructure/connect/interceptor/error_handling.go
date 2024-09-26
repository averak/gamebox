package interceptor

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/infrastructure/connect/error_response"
	"github.com/averak/gamebox/protobuf/api/api_errors"
)

func NewErrorHandlingInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			resp, err := next(ctx, req)
			if err != nil {
				var e error_response.Error
				if ok := errors.As(err, &e); ok {
					return nil, e.ConnectError()
				}

				var message string
				if config.Get().GetDebug() {
					message = fmt.Sprintf("unknown error occurred: %v", err)
				} else {
					message = "unknown error occurred"
				}
				return nil, error_response.New(api_errors.ErrorCode_COMMON_UNKNOWN, api_errors.ErrorSeverity_ERROR_SEVERITY_ERROR, message)
			}
			return resp, nil
		}
	}
}
