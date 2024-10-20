package model

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInsufficientCoins   = errors.New("insufficient coins")
	ErrCoinsMustBePositive = errors.New("coins must be positive")
)

type Wallet struct {
	UserID  uuid.UUID
	Balance Coins
}

func (w *Wallet) Deposit(coins Coins) error {
	balance, err := NewCoins(int(w.Balance + coins))
	if err != nil {
		return err
	}
	w.Balance = balance
	return nil
}

func (w *Wallet) Withdraw(coins Coins) error {
	if w.Balance < coins {
		return fmt.Errorf("%w: %d < %d", ErrInsufficientCoins, w.Balance, coins)
	}
	w.Balance -= coins
	return nil
}

// Coins は、ゲーム内通貨を表します。
type Coins int

func NewCoins(v int) (Coins, error) {
	if v < 0 {
		return 0, fmt.Errorf("%w: %d", ErrCoinsMustBePositive, v)
	}
	return Coins(v), nil
}

func (c Coins) IsZero() bool {
	return c == 0
}
