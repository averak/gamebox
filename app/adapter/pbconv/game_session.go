package pbconv

import (
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/protobuf/resource"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToGameSessionPb(sess model.GameSession) *resource.GameSession {
	var finishedAt *timestamppb.Timestamp
	if !sess.FinishedAt.IsZero() {
		finishedAt = timestamppb.New(sess.FinishedAt)
	}
	return &resource.GameSession{
		SessionId:  sess.ID.String(),
		GameId:     ToGameIDPb(sess.GameID),
		Status:     ToGameStatusPb(sess.Status),
		Result:     ToGameResultPb(sess.Result),
		Wager:      int64(sess.Wager),
		Payout:     int64(sess.Payout),
		StartedAt:  timestamppb.New(sess.StartedAt),
		FinishedAt: finishedAt,
	}
}

func ToGameIDPb(id model.GameID) resource.GameID {
	switch id {
	case model.GameIDSolitaire:
		return resource.GameID_GAME_ID_SOLITAIRE
	case model.GameIDBlackjack:
		return resource.GameID_GAME_ID_BLACKJACK
	case model.GameIDJanken:
		return resource.GameID_GAME_ID_JANKEN
	}
	return resource.GameID_GAME_ID_UNSPECIFIED
}

func ToGameStatusPb(status model.GameStatus) resource.GameStatus {
	switch status {
	case model.GameStatusPlaying:
		return resource.GameStatus_GAME_STATUS_PLAYING
	case model.GameStatusFinished:
		return resource.GameStatus_GAME_STATUS_FINISHED
	}
	return resource.GameStatus_GAME_STATUS_UNSPECIFIED
}

func ToGameResultPb(result model.GameResult) resource.GameResult {
	switch result {
	case model.GameResultWin:
		return resource.GameResult_GAME_RESULT_WIN
	case model.GameResultLose:
		return resource.GameResult_GAME_RESULT_LOSE
	case model.GameResultDraw:
		return resource.GameResult_GAME_RESULT_DRAW
	case model.GameResultUnknown:
		return resource.GameResult_GAME_RESULT_UNSPECIFIED
	}
	return resource.GameResult_GAME_RESULT_UNSPECIFIED
}
