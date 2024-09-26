package echo_handler

import (
	"context"

	"github.com/averak/gamebox/app/infrastructure/connect/advice"
	"github.com/averak/gamebox/app/usecase/echo_usecase"
	"github.com/averak/gamebox/protobuf/api/debug"
	"github.com/averak/gamebox/protobuf/api/debug/debugconnect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type handler struct {
	uc *echo_usecase.Usecase
}

func NewHandler(uc *echo_usecase.Usecase, advice advice.Advice) debugconnect.EchoServiceHandler {
	return debug.NewEchoServiceHandler(&handler{uc: uc}, advice)
}

func (h handler) EchoV1(ctx context.Context, req *advice.Request[*debug.EchoServiceEchoV1Request]) (*debug.EchoServiceEchoV1Response, error) {
	result, err := h.uc.Echo(ctx, req.GameContext(), req.Msg().GetMessage())
	if err != nil {
		return nil, err
	}
	return &debug.EchoServiceEchoV1Response{
		Message:   result.Message,
		Timestamp: timestamppb.New(result.Timestamp),
	}, nil
}
