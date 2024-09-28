package bdd

import (
	"testing"
)

type Testcase[GIVEN any, WHEN any, THEN any] struct {
	Name      string
	Given     GIVEN
	Behaviors []Behavior[WHEN, THEN]
}

func (tc Testcase[GIVEN, WHEN, THEN]) Run(t *testing.T, f func(t *testing.T, given GIVEN, when WHEN, then THEN)) {
	t.Run(tc.Name, func(t *testing.T) {
		for _, b := range tc.Behaviors {
			t.Run(b.Name, func(t *testing.T) {
				f(t, tc.Given, b.When, b.Then)
			})
		}
	})
}

type Behavior[WHEN any, THEN any] struct {
	Name string
	When WHEN
	Then THEN
}
