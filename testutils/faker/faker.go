package faker

import (
	"github.com/google/uuid"
)

var (
	space = uuid.New()
)

func UUIDv5(key string) uuid.UUID {
	return uuid.NewSHA1(space, []byte(key))
}
