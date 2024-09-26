package usecase

import (
	"github.com/averak/gamebox/app/usecase/echo_usecase"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	echo_usecase.NewUsecase,
)
