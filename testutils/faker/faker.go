package faker

import (
	"time"

	"github.com/google/uuid"
)

var (
	space = uuid.New()
)

func UUIDv5(key string) uuid.UUID {
	return uuid.NewSHA1(space, []byte(key))
}

func Email() string {
	return uuid.NewString() + "@example.com"
}

func MaxTime() time.Time {
	return time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
}
