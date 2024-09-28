package janken_handler

import (
	"context"

	"github.com/averak/gamebox/app/infrastructure/connect/advice"
	"github.com/averak/gamebox/protobuf/api"
	"github.com/averak/gamebox/protobuf/api/apiconnect"
)

type handler struct {
}

func NewHandler(advice advice.Advice) apiconnect.JankenServiceHandler {
	return api.NewJankenServiceHandler(&handler{}, advice)
}

func (h handler) ChooseHandV1(ctx context.Context, req *advice.Request[*api.JankenServiceChooseHandV1Request]) (*api.JankenServiceChooseHandV1Response, error) {
	// TODO: implement me
	return &api.JankenServiceChooseHandV1Response{}, nil
}

func (h handler) ChooseHandV1Errors(errs *api.JankenServiceChooseHandV1Errors) {
}
