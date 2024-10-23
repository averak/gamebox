package janken_handler

import (
	"context"

	"github.com/averak/gamebox/app/infrastructure/connect/aop"
	"github.com/averak/gamebox/protobuf/api"
	"github.com/averak/gamebox/protobuf/api/apiconnect"
)

type handler struct {
}

func NewHandler(proxy aop.Proxy) apiconnect.JankenServiceHandler {
	return api.NewJankenServiceHandler(&handler{}, proxy)
}

func (h handler) ChooseHandV1(ctx context.Context, req *aop.Request[*api.JankenServiceChooseHandV1Request]) (*api.JankenServiceChooseHandV1Response, error) {
	// TODO: implement me
	return &api.JankenServiceChooseHandV1Response{}, nil
}

func (h handler) ChooseHandV1Errors(errs *api.JankenServiceChooseHandV1Errors) {
}
