package main

import (
	"fmt"
	"strings"
	"time"
)

type Match struct {
	FirstTeam  string
	SecondTeam string
	Link       string
	Time       time.Time
	Stars      int
}

func NewMatch(firstTeam, secondTeam, link string, stars int, t time.Time) Match {
	return Match{
		FirstTeam:  firstTeam,
		SecondTeam: secondTeam,
		Link:       link,
		Time:       t,
		Stars:      stars,
	}
}

func (m *Match) IsToday() bool {
	return m.Time.Truncate(24 * time.Hour).Equal(time.Now().Truncate(24 * time.Hour))
}

func (m *Match) String() string {
	stars := strings.Repeat("*", m.Stars)
	time := fmt.Sprintf("%d:%d", m.Time.Hour(), m.Time.Minute())
	return fmt.Sprintf("%15s vs %-15s %5s %s", m.FirstTeam, m.SecondTeam, time, stars)
}
