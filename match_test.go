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
		Time:       time.Date(2018, 2, 3, 14, 30, 1, 0, time.UTC),
		Link:       "/link",
	}

	assert.Equal(t, "          first vs second          14:30 ***", m.String())
}
