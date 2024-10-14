package repoimpl

import (
	"github.com/averak/gamebox/app/adapter/repoimpl/echo_repoimpl"
	"github.com/averak/gamebox/app/adapter/repoimpl/game_session_repoimpl"
	"github.com/averak/gamebox/app/adapter/repoimpl/janken_session_repoimpl"
	"github.com/averak/gamebox/app/adapter/repoimpl/user_repoimpl"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	echo_repoimpl.NewRepository,
	game_session_repoimpl.NewRepository,
	janken_session_repoimpl.NewRepository,
	user_repoimpl.NewRepository,
)
