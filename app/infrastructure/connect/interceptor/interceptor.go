package interceptor

import (
	"connectrpc.com/connect"
	"github.com/averak/gamebox/app/core/build_info"
)

func New() []connect.Interceptor {
	// 上から順に実行される。
	// ただし、connect.UnaryFunc の戻り値に対して処理する interceptor は下から順番に実行される。
	return []connect.Interceptor{
		NewPanicRecoverInterceptor(),
		NewTraceInterceptor(),
		NewErrorHandlingInterceptor(),
		NewAccessLogInterceptor(),
		NewResponseMetadataInterceptor(build_info.ServerVersion()),
	}
}
