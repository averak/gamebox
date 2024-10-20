package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/averak/gamebox/app/core/logger"
	"github.com/google/uuid"
)

var (
	ErrGameAlreadyPlaying  = errors.New("game is already playing")
	ErrGameAlreadyFinished = errors.New("game is already finished")
	ErrGameNotPlaying      = errors.New("game is not playing")
	ErrGameWagerIsInvalid  = errors.New("game wager is invalid")
)

type GameID int

const (
	GameIDSolitaire GameID = iota + 1
	GameIDBlackjack
	GameIDJanken
)

type GameStatus int

const (
	GameStatusPlaying GameStatus = iota + 1
	GameStatusFinished
)

type GameResult int

const (
	GameResultUnknown GameResult = iota
	GameResultWin
	GameResultLose
	GameResultDraw
)

// GameSession は、PvE ゲームのセッションを表します。
type GameSession struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	GameID     GameID
	Status     GameStatus
	Result     GameResult
	Wager      Coins
	Payout     Coins
	StartedAt  time.Time
	FinishedAt time.Time
}

func NewGameSession(id uuid.UUID, userID uuid.UUID, gameID GameID, status GameStatus, result GameResult, wager Coins, payout Coins, startedAt time.Time, finishedAt time.Time) (GameSession, error) {
	if wager <= 0 {
		return GameSession{}, errors.New("wager must be positive")
	}
	return GameSession{
		ID:         id,
		UserID:     userID,
		GameID:     gameID,
		Status:     status,
		Result:     result,
		Wager:      wager,
		Payout:     payout,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
	}, nil
}

func (g *GameSession) FinishPlaying(result GameResult, now time.Time) (Coins, error) {
	if g.Status == GameStatusFinished {
		return 0, ErrGameAlreadyFinished
	}

	switch result {
	case GameResultWin:
		// 現時点では、賭け金の2倍を払い戻す。
		g.Payout = g.Wager * 2
	case GameResultLose:
		// do nothing
	case GameResultDraw:
		g.Payout = g.Wager
	case GameResultUnknown:
		return 0, errors.New("invalid game result")
	}

	g.Status = GameStatusFinished
	g.Result = result
	g.FinishedAt = now
	return g.Payout, nil
}

func (g GameSession) IsPlaying() bool {
	return g.Status == GameStatusPlaying
}

type GameSessionStartService struct {
	userID          uuid.UUID
	playingSessions []GameSession
}

func NewGameSessionService(ctx context.Context, userID uuid.UUID, playingSessions []GameSession) (GameSessionStartService, error) {
	var gameIDs = make(map[GameID]struct{})
	for _, session := range playingSessions {
		if !session.IsPlaying() {
			return GameSessionStartService{}, fmt.Errorf("invalid game session status: %v", session.Status)
		}

		if _, ok := gameIDs[session.GameID]; ok {
			// 何らかの不具合により同じゲームを同時進行した状態になると、進行不能になる可能性がある。
			// 進行不能バグを許容し、緊急対応により解決する方針とする。
			logger.Critical(ctx, map[string]any{
				"error":  "multiple game sessions are playing for the same game",
				"userID": session.UserID,
				"gameID": session.GameID,
			})
			return GameSessionStartService{}, errors.New("multiple game sessions are playing for the same game")
		}
		gameIDs[session.GameID] = struct{}{}
	}
	return GameSessionStartService{userID, playingSessions}, nil
}

func (s *GameSessionStartService) StartPlaying(id uuid.UUID, gameID GameID, wager Coins, now time.Time) (GameSession, error) {
	for _, session := range s.playingSessions {
		if session.GameID == gameID {
			return GameSession{}, ErrGameAlreadyPlaying
		}
	}
	if wager.IsZero() {
		return GameSession{}, ErrGameWagerIsInvalid
	}
	sess, err := NewGameSession(id, s.userID, gameID, GameStatusPlaying, GameResultUnknown, wager, 0, now, time.Time{})
	if err != nil {
		return GameSession{}, err
	}
	s.playingSessions = append(s.playingSessions, sess)
	return sess, nil
}
