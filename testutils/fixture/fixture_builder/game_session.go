package fixture_builder

import (
	"testing"
	"time"

	"github.com/averak/gamebox/app/domain/model"
	"github.com/google/uuid"
)

type GameSessionBuilder struct {
	data model.GameSession
}

func NewGameSessionBuilder(t *testing.T, userID uuid.UUID, gameSessionID uuid.UUID) *GameSessionBuilder {
	t.Helper()

	wager, err := model.NewCoins(100)
	if err != nil {
		t.Fatal(err)
	}
	payout, err := model.NewCoins(200)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	data, err := model.NewGameSession(gameSessionID, userID, model.GameIDSolitaire, model.GameStatusFinished, model.GameResultWin, wager, payout, now.Add(-time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	return &GameSessionBuilder{
		data: data,
	}
}

func (b GameSessionBuilder) Build() model.GameSession {
	return b.data
}

func (b *GameSessionBuilder) GameID(v model.GameID) *GameSessionBuilder {
	b.data.GameID = v
	return b
}

func (b *GameSessionBuilder) Status(v model.GameStatus) *GameSessionBuilder {
	b.data.Status = v
	return b
}

func (b *GameSessionBuilder) StartedAt(v time.Time) *GameSessionBuilder {
	b.data.StartedAt = v
	return b
}

func (b *GameSessionBuilder) FinishedAt(v time.Time) *GameSessionBuilder {
	b.data.FinishedAt = v
	return b
}
