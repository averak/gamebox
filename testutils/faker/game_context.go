package faker

import (
	"time"

	"github.com/averak/gamebox/app/core/game_context"
	"github.com/google/uuid"
)

type GameContextBuilder struct {
	data game_context.GameContext
}

func NewGameContextBuilder() *GameContextBuilder {
	return &GameContextBuilder{
		data: game_context.NewGameContext(uuid.New(), time.Now()),
	}
}

func (b *GameContextBuilder) Build() game_context.GameContext {
	return b.data
}

func (b *GameContextBuilder) IdempotencyKey(idempotencyKey uuid.UUID) *GameContextBuilder {
	b.data = game_context.NewGameContext(idempotencyKey, b.data.Now())
	return b
}

func (b *GameContextBuilder) Now(now time.Time) *GameContextBuilder {
	b.data = game_context.NewGameContext(b.data.IdempotencyKey(), now)
	return b
}
