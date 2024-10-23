package aop

import (
	"context"
	"errors"

	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/core/game_context"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/app/infrastructure/connect/error_response"
	"github.com/averak/gamebox/app/infrastructure/connect/mdval"
	"github.com/averak/gamebox/protobuf/custom_option"
	"google.golang.org/protobuf/proto"
)

type (
	MethodOption                   = custom_option.MethodOption
	MethodErrDefinition            = custom_option.MethodErrorDefinition
	Method[REQ, RES proto.Message] func(context.Context, *Request[REQ]) (RES, error)

	// Proxy は、rpc method の前後で cross-cutting concern を実行するための関数です。
	// interceptor だと共通化しづらい処理を、ここで実行します。
	Proxy func(context.Context, proto.Message, mdval.IncomingMD, *MethodInfo, func(context.Context, game_context.GameContext, *model.User, mdval.IncomingMD) (proto.Message, error)) error
)

func NewProxy(conn transaction.Connection, userRepo repository.UserRepository) Proxy {
	return func(ctx context.Context, req proto.Message, incomingMD mdval.IncomingMD, info *MethodInfo, method func(context.Context, game_context.GameContext, *model.User, mdval.IncomingMD) (proto.Message, error)) error {
		params, err := fixPreconditionParams(ctx, incomingMD)
		if err != nil {
			return err
		}
		gctx := params.GameContext()

		var principal *model.User
		if !info.Option().GetSkipAuthenticate() {
			principal, err = checkSession(ctx, config.Get(), userRepo, conn, incomingMD, gctx.Now())
			if err != nil {
				return err
			}

			if !info.Option().GetSkipUserStatusCheck() {
				err = checkUserStatus(*principal)
				if err != nil {
					return err
				}
			}
		}

		_, err = method(ctx, gctx, principal, incomingMD)
		if err != nil {
			if errDef, ok := info.FindErrorDefinition(err); ok {
				return error_response.New(errDef.GetCode(), errDef.GetSeverity(), errDef.GetMessage())
			}
			return err
		}
		return nil
	}
}

type MethodInfo struct {
	opt       *MethodOption
	errCauses map[error]*MethodErrDefinition
}

func NewMethodInfo(opt *MethodOption, errCauses map[error]*MethodErrDefinition) *MethodInfo {
	return &MethodInfo{
		opt:       opt,
		errCauses: errCauses,
	}
}

func (m *MethodInfo) Option() *MethodOption {
	return m.opt
}

func (m MethodInfo) FindErrorDefinition(err error) (*MethodErrDefinition, bool) {
	for cause, def := range m.errCauses {
		if errors.Is(err, cause) {
			return def, true
		}
	}
	return nil, false
}

type Request[T any] struct {
	msg       T
	gctx      game_context.GameContext
	principal *model.User
}

func NewRequest[T proto.Message](msg T, gctx game_context.GameContext, principal *model.User) *Request[T] {
	return &Request[T]{
		msg:       msg,
		gctx:      gctx,
		principal: principal,
	}
}

func (r Request[T]) Msg() T {
	return r.msg
}

func (r Request[T]) GameContext() game_context.GameContext {
	return r.gctx
}

// Principal は、認証されたユーザーを返します。
// 認証必須の API は必ず true を返すので、わざわざ戻り値をチェックする必要はありません。
func (r Request[T]) Principal() (model.User, bool) {
	if r.principal == nil {
		return model.User{}, false
	}
	return *r.principal, true
}
