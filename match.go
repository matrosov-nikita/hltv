package main

import (
	"fmt"
	"strings"
	"time"
)

// Match represents math from hltv.org site.
type Match struct {
	FirstTeam  string
	SecondTeam string
	Link       string
	Time       time.Time
	Stars      int
}

const (
	day                 = 24 * time.Hour
	maxTeamLength       = 18
	trailingDotsForTeam = "..."
)

func cutLongTeamName(team string) string {
	if len(team) > maxTeamLength {
		return team[:maxTeamLength-len(trailingDotsForTeam)] + trailingDotsForTeam
	}

	return team
}

// NewMatch creates new match from given info.
func NewMatch(firstTeam, secondTeam, link string, stars int, t time.Time) Match {
	return Match{
		FirstTeam:  firstTeam,
		SecondTeam: secondTeam,
		Link:       link,
		Time:       t,
		Stars:      stars,
	}
}

// IsToday verifies that match will be today.
func (m *Match) IsToday() bool {
	return m.Time.Truncate(day).Equal(time.Now().Truncate(day))
}

func (m *Match) String() string {
	stars := strings.Repeat("*", m.Stars)
	time := fmt.Sprintf("%d:%.2d", m.Time.Hour(), m.Time.Minute())

	return fmt.Sprintf("%18s vs %-18s %5s %s",
		cutLongTeamName(m.FirstTeam), cutLongTeamName(m.SecondTeam), time, stars)
}
