package interceptor

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/app/infrastructure/connect/error_response"
	"github.com/averak/gamebox/protobuf/api/api_errors"
	"github.com/averak/gamebox/testutils/testconnect"
	"github.com/stretchr/testify/require"
)

func TestErrorHandlingInterceptor(t *testing.T) {
	type args struct {
		next connect.UnaryFunc
	}
	tests := []struct {
		name string
		args args
		then func(*testing.T, error)
	}{
		{
			name: "発生したエラーが error_response.Error の場合 => そのエラーを返す",
			args: args{
				next: func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
					return nil, error_response.New(api_errors.ErrorCode_COMMON_INVALID_SESSION, api_errors.ErrorSeverity_ERROR_SEVERITY_ERROR, "error")
				},
			},
			then: func(t *testing.T, err error) {
				testconnect.AssertErrorCode(t, api_errors.ErrorCode_COMMON_INVALID_SESSION, err)
			},
		},
		{
			name: "発生したエラーが error_response.Error ではない => UNKNOWN エラーを返す",
			args: args{
				next: func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
					return nil, errors.New("error")
				},
			},
			then: func(t *testing.T, err error) {
				testconnect.AssertErrorCode(t, api_errors.ErrorCode_COMMON_UNKNOWN, err)
			},
		},
		{
			name: "エラーが発生しない => エラーハンドリングを行わない",
			args: args{
				next: func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
					return nil, nil
				},
			},
			then: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := NewErrorHandlingInterceptor()
			_, err := interceptor(tt.args.next)(context.Background(), connect.NewRequest[any](nil))
			tt.then(t, err)
		})
	}
}
