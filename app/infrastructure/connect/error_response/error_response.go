package error_response

import (
	"errors"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/protobuf/api/api_errors"
	rpccode "google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type CodeType interface {
	api_errors.ErrorCode_Common | api_errors.ErrorCode_Method

	Descriptor() protoreflect.EnumDescriptor
	Number() protoreflect.EnumNumber
}

// New は Connect RPC のエラーを生成します。
func New[T CodeType](code T, severity api_errors.ErrorSeverity, message string) Error {
	var (
		grpcCode     rpccode.Code
		handlingType api_errors.ErrorHandlingType
	)
	values := code.Descriptor().Values()
	for i := range values.Len() {
		v := values.Get(i)
		if v.Number() == code.Number() {
			if c, ok := proto.GetExtension(v.Options(), api_errors.E_GrpcCode).(rpccode.Code); ok {
				grpcCode = c
			}
			if h, ok := proto.GetExtension(v.Options(), api_errors.E_ErrorHandlingType).(api_errors.ErrorHandlingType); ok {
				handlingType = h
			}
			break
		}
	}

	detail := &api_errors.ErrorDetail{
		ErrorCode:         int64(code.Number()),
		ErrorHandlingType: handlingType,
	}
	return newError(connect.Code(grpcCode), severity, detail, message)
}

type Error struct {
	connectErr *connect.Error
	severity   api_errors.ErrorSeverity
}

func newError(code connect.Code, severity api_errors.ErrorSeverity, detail *api_errors.ErrorDetail, message string) Error {
	ce := connect.NewError(code, errors.New(message))
	if d, err := connect.NewErrorDetail(detail); err == nil {
		ce.AddDetail(d)
	}
	return Error{ce, severity}
}

func (e Error) Error() string {
	return e.connectErr.Error()
}

func (e Error) Severity() api_errors.ErrorSeverity {
	return e.severity
}

func (e Error) ConnectError() *connect.Error {
	return e.connectErr
}
