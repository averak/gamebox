package repository

import (
	"context"
	"errors"

	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	Get(ctx context.Context, tx transaction.Transaction, userID uuid.UUID) (model.User, error)
	Save(ctx context.Context, tx transaction.Transaction, user model.User) error
}
