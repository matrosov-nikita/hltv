package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMatchFormat(t *testing.T) {
	m := &Match{
		FirstTeam:  "first",
		SecondTeam: "second",
		Stars:      3,
		Time:       time.Date(2018, 2, 3, 5, 7, 1, 0, time.UTC),
		Link:       "/link",
	}

	assert.Equal(t, "             first vs second              5:07 ***", m.String())
}

func TestMatchCutLongTeamName(t *testing.T) {
	m := NewMatch("Pretty long name for first team", "Very long name for second team",
		"/link", 3, time.Date(2018, 2, 3, 5, 7, 1, 0, time.UTC))

	assert.Equal(t, "Pretty long nam... vs Very long name ...  5:07 ***", m.String())
}
