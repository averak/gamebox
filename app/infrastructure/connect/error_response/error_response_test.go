package error_response_test

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/app/infrastructure/connect/error_response"
	"github.com/averak/gamebox/protobuf/api/api_errors"
	"github.com/averak/gamebox/testutils/testconnect"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args[CODE error_response.CodeType] struct {
		code     CODE
		severity api_errors.ErrorSeverity
		message  string
	}
	type testCase[CODE error_response.CodeType] struct {
		name string
		args args[CODE]
		then func(t *testing.T, got error_response.Error)
	}
	tests := []testCase[api_errors.ErrorCode_Common]{
		{
			name: "[ケース1] COMMON_UNKNOWN",
			args: args[api_errors.ErrorCode_Common]{
				code: api_errors.ErrorCode_COMMON_UNKNOWN,
			},
			then: func(t *testing.T, got error_response.Error) {
				assert.Equal(t, connect.CodeUnknown, got.ConnectError().Code())

				wantDetail := &api_errors.ErrorDetail{
					ErrorCode:         int64(api_errors.ErrorCode_COMMON_UNKNOWN),
					ErrorHandlingType: api_errors.ErrorHandlingType_ERROR_HANDLING_TYPE_TEMPORARY,
				}
				assert.EqualExportedValues(t, wantDetail, testconnect.GetErrorDetail(got.ConnectError()))
			},
		},
		{
			name: "[ケース2] METHOD_COMMON_INVALID_SESSION",
			args: args[api_errors.ErrorCode_Common]{
				code: api_errors.ErrorCode_COMMON_INVALID_SESSION,
			},
			then: func(t *testing.T, got error_response.Error) {
				assert.Equal(t, connect.CodeUnauthenticated, got.ConnectError().Code())

				wantDetail := &api_errors.ErrorDetail{
					ErrorCode:         int64(api_errors.ErrorCode_COMMON_INVALID_SESSION),
					ErrorHandlingType: api_errors.ErrorHandlingType_ERROR_HANDLING_TYPE_RECOVERABLE,
				}
				assert.EqualExportedValues(t, wantDetail, testconnect.GetErrorDetail(got.ConnectError()))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := error_response.New(tt.args.code, tt.args.severity, tt.args.message)
			tt.then(t, err)
		})
	}
}
