package repository

import (
	"context"
	"errors"

	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/google/uuid"
)

var (
	ErrJankenSessionNotFound = errors.New("janken session not found")
)

type JankenSessionRepository interface {
	Get(ctx context.Context, tx transaction.Transaction, gameSessionID uuid.UUID) (model.JankenSession, error)
	Save(ctx context.Context, tx transaction.Transaction, session model.JankenSession) error
}
