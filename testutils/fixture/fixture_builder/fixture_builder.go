package fixture_builder

import (
	"github.com/averak/gamebox/app/domain/model"
)

type Fixture struct {
	Users        []model.User
	GameSessions []model.GameSession
}

type FixtureBuilder struct {
	data *Fixture
}

func New() *FixtureBuilder {
	return &FixtureBuilder{
		data: &Fixture{},
	}
}

func (b FixtureBuilder) Build() Fixture {
	return *b.data
}

func (b *FixtureBuilder) User(v ...model.User) *FixtureBuilder {
	b.data.Users = append(b.data.Users, v...)
	return b
}

func (b *FixtureBuilder) GameSession(v ...model.GameSession) *FixtureBuilder {
	b.data.GameSessions = append(b.data.GameSessions, v...)
	return b
}
