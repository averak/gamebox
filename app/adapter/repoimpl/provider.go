package repoimpl

import (
	"github.com/averak/gamebox/app/adapter/repoimpl/echo_repoimpl"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	echo_repoimpl.NewRepository,
)
