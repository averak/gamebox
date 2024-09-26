package handler

import (
	"net/http"

	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	New,
)

func New() *http.ServeMux {
	return http.NewServeMux()
}
