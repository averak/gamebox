package connect

import (
	"context"
	"net/http"

	"github.com/averak/gamebox/app/core/game_context"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/infrastructure/connect/advice"
	"github.com/averak/gamebox/app/infrastructure/connect/mdval"
	"google.golang.org/protobuf/proto"
)

func Invoke[REQ, RES proto.Message](ctx context.Context, req REQ, header http.Header, info *advice.MethodInfo, method func(context.Context, *advice.Request[REQ]) (RES, error), adv advice.Advice) (RES, error) {
	var res RES
	wrap := func(ctx context.Context, gctx game_context.GameContext, principal *model.User, incomingMD mdval.IncomingMD) (proto.Message, error) {
		var err error
		res, err = method(ctx, advice.NewRequest(req, gctx, principal))
		return res, err
	}
	return res, adv(ctx, req, mdval.NewIncomingMD(header), info, wrap)
}
