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
)

type GameID int

const (
	GameIDSolitaire GameID = iota + 1
	GameIDBlackjack
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
)

type GameSession struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	GameID     GameID
	Status     GameStatus
	Result     GameResult
	Wager      int
	Payout     int
	StartedAt  time.Time
	FinishedAt time.Time
}

func NewGameSession(id uuid.UUID, userID uuid.UUID, gameID GameID, wager int, payout int, status GameStatus, result GameResult, startedAt time.Time, finishedAt time.Time) GameSession {
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
	}
}

func (s *GameSession) FinishPlaying(result GameResult, wallet Wallet, now time.Time) (GameSession, Wallet, error) {
	if s.Status == GameStatusFinished {
		return GameSession{}, Wallet{}, ErrGameAlreadyFinished
	}

	s.Status = GameStatusFinished
	s.Result = result
	s.FinishedAt = now

	switch result {
	case GameResultWin:
		// 後続処理で報酬を獲得させるため、ここでは何もしない。
	case GameResultLose:
		return *s, wallet, nil
	case GameResultUnknown:
		return GameSession{}, Wallet{}, errors.New("invalid game result")
	}

	s.Payout = s.Wager * 2
	if err := wallet.deposit(s.Payout); err != nil {
		return GameSession{}, Wallet{}, err
	}
	return *s, wallet, nil
}

type GameSessionService struct {
	userID          uuid.UUID
	playingSessions []GameSession
}

func NewGameSessionService(ctx context.Context, userID uuid.UUID, playingSessions []GameSession) (GameSessionService, error) {
	var gameIDs = make(map[GameID]struct{})
	for _, session := range playingSessions {
		if session.Status != GameStatusPlaying {
			return GameSessionService{}, fmt.Errorf("invalid game session status: %v", session.Status)
		}

		if _, ok := gameIDs[session.GameID]; ok {
			// 何らかの不具合により同じゲームを同時進行した状態になると、進行不能になる可能性がある。
			// 進行不能バグを許容し、緊急対応により解決する方針とする。
			logger.Critical(ctx, map[string]any{
				"error":  "multiple game sessions are playing for the same game",
				"userID": session.UserID,
				"gameID": session.GameID,
			})
			return GameSessionService{}, errors.New("multiple game sessions are playing for the same game")
		}
		gameIDs[session.GameID] = struct{}{}
	}
	return GameSessionService{userID, playingSessions}, nil
}

func (s *GameSessionService) StartPlaying(id uuid.UUID, gameID GameID, wager int, now time.Time) (GameSession, error) {
	for _, session := range s.playingSessions {
		if session.GameID == gameID {
			return GameSession{}, ErrGameAlreadyPlaying
		}
	}
	sess := NewGameSession(id, s.userID, gameID, wager, 0, GameStatusPlaying, GameResultUnknown, now, time.Time{})
	s.playingSessions = append(s.playingSessions, sess)
	return sess, nil
}
