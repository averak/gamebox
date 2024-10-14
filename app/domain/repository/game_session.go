package repository

import (
	"context"
	"errors"

	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/google/uuid"
)

var (
	ErrGameSessionNotFound = errors.New("game session not found")
)

type GameSessionRepository interface {
	Get(ctx context.Context, tx transaction.Transaction, userID uuid.UUID, gameSessionID uuid.UUID) (model.GameSession, error)
	Save(ctx context.Context, tx transaction.Transaction, sessions ...model.GameSession) error
}
