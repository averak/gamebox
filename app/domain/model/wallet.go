package model

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type Wallet struct {
	UserID  uuid.UUID
	Balance int
}

func (w *Wallet) deposit(amount int) error {
	if amount < 0 {
		return errors.New("amount must be positive")
	}
	w.Balance += amount
	return nil
}

func (w *Wallet) withdraw(amount int) error {
	if amount < 0 {
		return errors.New("amount must be positive")
	}
	if w.Balance < amount {
		return fmt.Errorf("%w: %d < %d", ErrInsufficientBalance, w.Balance, amount)
	}
	w.Balance -= amount
	return nil
}
