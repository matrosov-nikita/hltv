package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type MatchesFilterSuite struct {
	suite.Suite

	filterStub *MatchesFilterStub
}

func (s *MatchesFilterSuite) SetupTest() {
	s.filterStub = &MatchesFilterStub{}
}

func (s MatchesFilterSuite) TestFilterReturnsTodayMatches() {
	matches := NewFilterTodayMatches(s.filterStub).getMatches()
	s.Len(matches, 3)
}

func (s MatchesFilterSuite) TestFilterReturnsStarredMatches() {
	matches := NewStarredFilter(s.filterStub).getMatches()
	s.Len(matches, 3)
}

func (s MatchesFilterSuite) TestFilterReturnsMatchesByTeam() {
	matches := NewTeamFilter(s.filterStub, "faze").getMatches()
	s.Len(matches, 1)
	s.Equal(matches[0].FirstTeam, "FaZe")
}

func TestMatchesFilterSuite(t *testing.T) {
	suite.Run(t, new(MatchesFilterSuite))
}

type MatchesFilterStub struct{}

func (f *MatchesFilterStub) getMatches() []Match {
	today := time.Now()
	yesterday := time.Now().Add(-time.Hour * 24)

	return []Match{
		NewMatch("first", "second", "/link", 0, today),
		NewMatch("third", "fourth", "/link", 5, yesterday),
		NewMatch("FaZe", "Navi", "/link", 1, today),
		NewMatch("first", "second", "/link", 5, yesterday),
		NewMatch("first", "second", "/link", 0, today),
	}
}
