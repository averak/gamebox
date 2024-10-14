package fixture_builder

import (
	"github.com/averak/gamebox/app/domain/model"
	"github.com/google/uuid"
)

type UserBuilder struct {
	data model.User
}

func NewUserBuilder(id uuid.UUID) *UserBuilder {
	return &UserBuilder{
		data: model.NewUser(id, model.UserStatusActive),
	}
}

func (b UserBuilder) Build() model.User {
	return b.data
}
