package game_handler

import (
	"context"

	"github.com/averak/gamebox/app/adapter/pbconv"
	"github.com/averak/gamebox/app/domain/repository"
	"github.com/averak/gamebox/app/infrastructure/connect/advice"
	"github.com/averak/gamebox/app/usecase/game_usecase"
	"github.com/averak/gamebox/protobuf/api"
	"github.com/averak/gamebox/protobuf/api/apiconnect"
	"github.com/google/uuid"
)

type handler struct {
	uc *game_usecase.Usecase
}

func NewHandler(uc *game_usecase.Usecase, advice advice.Advice) apiconnect.GameServiceHandler {
	return api.NewGameServiceHandler(&handler{uc}, advice)
}

func (h handler) GetSessionV1(ctx context.Context, req *advice.Request[*api.GameServiceGetSessionV1Request]) (*api.GameServiceGetSessionV1Response, error) {
	principal, _ := req.Principal()
	result, err := h.uc.GetSession(ctx, principal, uuid.MustParse(req.Msg().GetSessionId()))
	if err != nil {
		return nil, err
	}
	return &api.GameServiceGetSessionV1Response{
		Session: pbconv.ToGameSessionPb(result),
	}, nil
}

func (h handler) GetSessionV1Errors(errs *api.GameServiceGetSessionV1Errors) {
	errs.Map(repository.ErrGameSessionNotFound, errs.RESOURCE_NOT_FOUND)
}

func (h handler) ListPlayingSessionsV1(ctx context.Context, req *advice.Request[*api.GameServiceListPlayingSessionsV1Request]) (*api.GameServiceListPlayingSessionsV1Response, error) {
	// TODO: implement me
	return &api.GameServiceListPlayingSessionsV1Response{}, nil
}

func (h handler) StartPlayingV1(ctx context.Context, req *advice.Request[*api.GameServiceStartPlayingV1Request]) (*api.GameServiceStartPlayingV1Response, error) {
	// TODO: implement me
	return &api.GameServiceStartPlayingV1Response{}, nil
}

func (h handler) StartPlayingV1Errors(errs *api.GameServiceStartPlayingV1Errors) {
}
