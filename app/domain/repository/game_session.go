package repository

import (
	"context"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/google/uuid"
)

type GameSessionRepository interface {
	Save(ctx context.Context, tx transaction.Transaction, sessions ...model.GameSession) error
	Get(ctx context.Context, tx transaction.Transaction, userID uuid.UUID, status model.GameStatus) ([]model.GameSession, error)
}
