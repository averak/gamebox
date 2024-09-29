package model

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserDeactivated = errors.New("user is deactivated")
)

type UserStatus int

const (
	// 設定が未完了で、他ユーザには公開されていない状態。
	// サインイン後、設定を完了するまでこの状態が続く。
	UserStatusPending UserStatus = 1 + iota
	// アカウントが利用可能で、他ユーザにも公開されている状態。
	UserStatusActive
	// アカウントが削除され、利用不可な状態。
	// この状態になると、二度と利用できなくなる。
	UserStatusDeactivated
)

type User struct {
	ID     uuid.UUID
	Status UserStatus
}

func NewUser(id uuid.UUID, status UserStatus) User {
	return User{
		ID:     id,
		Status: status,
	}
}
