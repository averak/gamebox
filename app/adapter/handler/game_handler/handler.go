package game_handler

import (
	"context"

	"github.com/averak/gamebox/app/adapter/pbconv"
	"github.com/averak/gamebox/app/domain/model"
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
	principal, _ := req.Principal()
	result, err := h.uc.ListPlayingSession(ctx, principal)
	if err != nil {
		return nil, err
	}
	return &api.GameServiceListPlayingSessionsV1Response{
		Sessions: pbconv.ToGameSessionPbs(result),
	}, nil
}

func (h handler) StartPlayingV1(ctx context.Context, req *advice.Request[*api.GameServiceStartPlayingV1Request]) (*api.GameServiceStartPlayingV1Response, error) {
	principal, _ := req.Principal()
	gameID, err := pbconv.ToGameID(req.Msg().GetGameId())
	if err != nil {
		return nil, err
	}
	wager, err := model.NewCoins(int(req.Msg().GetWager()))
	if err != nil {
		return nil, err
	}
	result, err := h.uc.StartPlaying(ctx, req.GameContext(), principal, gameID, wager)
	if err != nil {
		return nil, err
	}
	return &api.GameServiceStartPlayingV1Response{
		Session: pbconv.ToGameSessionPb(result),
	}, nil
}

func (h handler) StartPlayingV1Errors(errs *api.GameServiceStartPlayingV1Errors) {
	errs.Map(model.ErrCoinsMustBePositive, errs.ILLEGAL_ARGUMENT)
	errs.Map(model.ErrGameWagerIsInvalid, errs.ILLEGAL_ARGUMENT)
	errs.Map(model.ErrGameAlreadyPlaying, errs.RESOURCE_CONFLICT)
	errs.Map(pbconv.ErrGameIDNotExists, errs.RESOURCE_NOT_FOUND)
}
