package game_session_handler

import (
	"context"

	"github.com/averak/gamebox/app/infrastructure/connect/advice"
	"github.com/averak/gamebox/protobuf/api"
	"github.com/averak/gamebox/protobuf/api/apiconnect"
)

type handler struct {
}

func NewHandler(advice advice.Advice) apiconnect.GameServiceHandler {
	return api.NewGameServiceHandler(&handler{}, advice)
}

func (h handler) GetSessionV1(ctx context.Context, req *advice.Request[*api.GameServiceGetSessionV1Request]) (*api.GameServiceGetSessionV1Response, error) {
	// TODO: implement me
	return &api.GameServiceGetSessionV1Response{}, nil
}

func (h handler) GetSessionV1Errors(errs *api.GameServiceGetSessionV1Errors) {
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
