package game_context

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

var (
	_ json.Marshaler = GameContext{}
)

// GameContext は、機能によらずアプリケーション横断的なコンテキストを提供します。
type GameContext struct {
	idempotencyKey uuid.UUID
	now            time.Time
}

func NewGameContext(idempotencyKey uuid.UUID, now time.Time) GameContext {
	return GameContext{
		idempotencyKey: idempotencyKey,
		now:            now,
	}
}

// IdempotencyKey はトランザクションの冪等性を保証するために利用される、一意な識別子です。
// https://developer.mozilla.org/ja/docs/Glossary/Idempotent
func (c GameContext) IdempotencyKey() uuid.UUID {
	return c.idempotencyKey
}

func (c GameContext) Now() time.Time {
	return c.now
}

func (c GameContext) JSON() map[string]interface{} {
	return map[string]interface{}{
		"idempotencyKey": c.idempotencyKey,
		"now":            c.now,
	}
}

func (c GameContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IdempotencyKey string `json:"idempotencyKey"`
		Now            string `json:"now"`
	}{
		IdempotencyKey: c.idempotencyKey.String(),
		Now:            c.now.Format(time.RFC3339Nano),
	})
}
