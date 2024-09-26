package interceptor

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/protobuf/api/api_errors"
	"github.com/averak/gamebox/testutils/testconnect"
	"github.com/stretchr/testify/assert"
)

func TestPanicRecoverInterceptor(t *testing.T) {
	type args struct {
		next connect.UnaryFunc
	}
	tests := []struct {
		name string
		args args
		then func(*testing.T, error)
	}{
		{
			name: "panicが発生しない場合 => 何もしない",
			args: args{
				next: func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
					return nil, nil
				},
			},
			then: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "panicが発生した場合 => エラーを返す",
			args: args{
				next: func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
					panic("panic")
				},
			},
			then: func(t *testing.T, err error) {
				testconnect.AssertErrorCode(t, api_errors.ErrorCode_COMMON_UNKNOWN, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := NewPanicRecoverInterceptor()
			_, err := interceptor(tt.args.next)(context.Background(), connect.NewRequest[any](nil))
			tt.then(t, err)
		})
	}
}
