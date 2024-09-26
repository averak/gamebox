package interceptor

import (
	"context"
	"errors"
	"testing"

	"github.com/averak/gamebox/app/infrastructure/connect/error_response"
	"github.com/averak/gamebox/protobuf/api/api_errors"
	"github.com/stretchr/testify/assert"
)

func Test_getSeverity(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want api_errors.ErrorSeverity
	}{
		{
			name: "error_response.Error の場合 => severity を取得できる",
			args: args{
				err: error_response.New(api_errors.ErrorCode_COMMON_UNKNOWN, api_errors.ErrorSeverity_ERROR_SEVERITY_CRITICAL, "unknown error occurred"),
			},
			want: api_errors.ErrorSeverity_ERROR_SEVERITY_CRITICAL,
		},
		{
			name: "error_response.Error 以外の場合 => ERROR",
			args: args{
				err: errors.New("error"),
			},
			want: api_errors.ErrorSeverity_ERROR_SEVERITY_ERROR,
		},
		{
			name: "context.Canceled の場合 => WARNING",
			args: args{
				err: context.Canceled,
			},
			want: api_errors.ErrorSeverity_ERROR_SEVERITY_WARNING,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getSeverity(tt.args.err), "getSeverity(%v)", tt.args.err)
		})
	}
}
