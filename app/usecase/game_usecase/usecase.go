package game_usecase

import (
	"github.com/averak/gamebox/app/domain/repository"
	"github.com/averak/gamebox/app/domain/repository/transaction"
)

type Usecase struct {
	conn            transaction.Connection
	gameSessionRepo repository.GameSessionRepository
}

func NewUsecase(conn transaction.Connection, gameSessionRepo repository.GameSessionRepository) *Usecase {
	return &Usecase{
		conn:            conn,
		gameSessionRepo: gameSessionRepo,
	}
}
