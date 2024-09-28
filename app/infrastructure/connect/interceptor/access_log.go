package interceptor

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/core/logger"
	"github.com/averak/gamebox/app/infrastructure/connect/error_response"
	"github.com/averak/gamebox/app/infrastructure/connect/mdval"
	"github.com/averak/gamebox/protobuf/api/api_errors"
)

type (
	hintKey struct{}
	hint    = map[string]any
)

func NewAccessLogInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			begin := time.Now()
			ctx = context.WithValue(ctx, hintKey{}, make(hint))

			resp, err := next(ctx, req)

			var request any
			if req != nil && !config.Get().GetLogging().GetOmitRequestBody() {
				request = req.Any()
			}
			var response any
			if resp != nil &&
				!config.Get().GetLogging().GetOmitResponseBody() {
				response = resp.Any()
			}
			payload := map[string]interface{}{
				"procedure":     req.Spec().Procedure,
				"requestedAt":   time.Now().UTC().Format(time.RFC3339Nano),
				"elapsedTimeMs": time.Since(begin).Milliseconds(),
				"request":       request,
				"response":      response,
				"incomingMD":    mdval.NewIncomingMD(req.Header()),
			}
			hnt, ok := getLogHint(ctx)
			if ok {
				payload["hint"] = hnt
			}

			if err == nil {
				logger.Info(ctx, payload)
			} else {
				payload["error"] = err.Error()
				switch getSeverity(err) {
				case api_errors.ErrorSeverity_ERROR_SEVERITY_UNSPECIFIED:
					// API スキーマで severity の設定漏れで UNSPECIFIED になることがあるが、その場合は ERROR として扱う。
					logger.Error(ctx, payload)

					// severity の設定漏れでクライアントにエラーを返すわけにはいかないが、不備を検知するためのログは出しておく。
					logger.Error(ctx, map[string]any{
						"error":     "severity is not specified",
						"procedure": req.Spec().Procedure,
					})
				case api_errors.ErrorSeverity_ERROR_SEVERITY_DEBUG:
					logger.Debug(ctx, payload)
				case api_errors.ErrorSeverity_ERROR_SEVERITY_INFO:
					logger.Info(ctx, payload)
				case api_errors.ErrorSeverity_ERROR_SEVERITY_NOTICE:
					logger.Notice(ctx, payload)
				case api_errors.ErrorSeverity_ERROR_SEVERITY_WARNING:
					logger.Warning(ctx, payload)
				case api_errors.ErrorSeverity_ERROR_SEVERITY_ERROR:
					logger.Error(ctx, payload)
				case api_errors.ErrorSeverity_ERROR_SEVERITY_CRITICAL:
					logger.Critical(ctx, payload)
				case api_errors.ErrorSeverity_ERROR_SEVERITY_ALERT:
					logger.Alert(ctx, payload)
				case api_errors.ErrorSeverity_ERROR_SEVERITY_EMERGENCY:
					logger.Emergency(ctx, payload)
				}
			}
			return resp, err
		}
	}
}

func getSeverity(err error) api_errors.ErrorSeverity {
	var e error_response.Error
	if errors.As(err, &e) {
		return e.Severity()
	} else if errors.Is(err, context.Canceled) {
		// クライアントが切断した場合は Warning ログを出す。
		return api_errors.ErrorSeverity_ERROR_SEVERITY_WARNING
	}
	return api_errors.ErrorSeverity_ERROR_SEVERITY_ERROR
}

func getLogHint(ctx context.Context) (hint, bool) {
	v, ok := ctx.Value(hintKey{}).(hint)
	return v, ok
}

func AddLogHint(ctx context.Context, key string, value any) {
	if h, ok := getLogHint(ctx); ok {
		h[key] = value
	}
}
