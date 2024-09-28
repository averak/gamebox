package model

import (
	"context"
	"math/rand"
	"time"

	"github.com/averak/gamebox/pkg/vector"
	"github.com/google/uuid"
)

type JankenHand int

const (
	JankenHandRock JankenHand = iota + 1
	JankenHandScissors
	JankenHandPaper
)

func (mine JankenHand) Battle(opponent JankenHand) GameResult {
	if mine == opponent {
		return GameResultDraw
	}
	if (mine == JankenHandRock && opponent == JankenHandScissors) ||
		(mine == JankenHandScissors && opponent == JankenHandPaper) ||
		(mine == JankenHandPaper && opponent == JankenHandRock) {
		return GameResultWin
	}
	return GameResultLose
}

// JankenSession は、PvE のジャンケンを表します。
type JankenSession struct {
	GameSessionID uuid.UUID
	Histories     []JankenHistory
}

func NewJankenSession(gameSessionID uuid.UUID) JankenSession {
	return JankenSession{
		GameSessionID: gameSessionID,
	}
}

func (j *JankenSession) Choose(session GameSession, hand JankenHand, now time.Time) (GameSession, Coins, JankenHistory, error) {
	if j.GameSessionID != session.ID || !session.IsPlaying() {
		return GameSession{}, 0, JankenHistory{}, ErrGameNotPlaying
	}

	oppHand := j.nextOpponentHand()
	history := NewJankenHistory(j.turn(), hand, oppHand)
	j.Histories = append(j.Histories, history)

	result := hand.Battle(oppHand)
	if result == GameResultDraw {
		return session, 0, history, nil
	}
	coins, err := session.FinishPlaying(result, now)
	if err != nil {
		return GameSession{}, 0, JankenHistory{}, err
	}
	return session, coins, history, nil
}

func (j JankenSession) turn() int {
	if len(j.Histories) == 0 {
		return 1
	}

	histories := vector.New(j.Histories).Sort(func(a, b JankenHistory) bool {
		return a.Turn < b.Turn
	})
	latest := histories[len(histories)-1]
	return latest.Turn + 1
}

// nextOpponentHand は、次のターンでの相手の手を返します。
// (ゲームセッションID, ターン数) の組み合わせで、決定論的に相手の手を決定します。
func (j JankenSession) nextOpponentHand() JankenHand {
	// ゲームセッションIDをシード値とする。
	var seed int64
	for _, b := range j.GameSessionID {
		seed += int64(b)
	}

	source := rand.NewSource(seed + int64(j.turn()))
	r := rand.New(source)
	return JankenHand(r.Intn(3) + 1)
}

type JankenHistory struct {
	Turn         int
	MyHand       JankenHand
	OpponentHand JankenHand
}

func NewJankenHistory(turn int, myHand JankenHand, opponentHand JankenHand) JankenHistory {
	return JankenHistory{
		Turn:         turn,
		MyHand:       myHand,
		OpponentHand: opponentHand,
	}
}

type JankenSessionStartService struct {
	GameSessionStartService
}

func NewJankenSessionStartService(ctx context.Context, userID uuid.UUID, playingSessions []GameSession) (JankenSessionStartService, error) {
	service, err := NewGameSessionService(ctx, userID, playingSessions)
	if err != nil {
		return JankenSessionStartService{}, err
	}
	return JankenSessionStartService{
		GameSessionStartService: service,
	}, nil
}

func (s JankenSessionStartService) StartPlaying(id uuid.UUID, wager Coins, now time.Time) (GameSession, JankenSession, error) {
	session, err := s.GameSessionStartService.StartPlaying(id, GameIDJanken, wager, now)
	if err != nil {
		return GameSession{}, JankenSession{}, err
	}
	return session, NewJankenSession(session.ID), nil
}
