package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type MatchesFilterSuite struct {
	suite.Suite

	todayMatch         Match
	yesterdayMatch     Match
	starredMatch       Match
	takeEveryMatchStub *TakeEveryMatchFilter
}

func (s *MatchesFilterSuite) SetupTest() {
	s.takeEveryMatchStub = &TakeEveryMatchFilter{}
	s.todayMatch = NewMatch("first", "second", "/link", 0, time.Now())
	s.yesterdayMatch = NewMatch("first", "second", "/link", 0, time.Now().Add(-24*time.Hour))
	s.starredMatch = NewMatch("first", "second", "/link", 5, time.Now())
}

func (s MatchesFilterSuite) TestTodayFilterTakeTodayMatches() {
	s.True(NewFilterTodayMatches(s.takeEveryMatchStub).TakeMatch(s.todayMatch))
}

func (s MatchesFilterSuite) TestTodayFilterSkipNotTodayMatches() {
	s.False(NewFilterTodayMatches(s.takeEveryMatchStub).TakeMatch(s.yesterdayMatch))
}

func (s MatchesFilterSuite) TestStarredFilterSkipUnstarredMatches() {
	unstarredMatch := s.todayMatch
	s.False(NewStarredFilter(s.takeEveryMatchStub).TakeMatch(unstarredMatch))
}

func (s MatchesFilterSuite) TestStarredFilterTakeAllStarredMatches() {
	s.True(NewStarredFilter(s.takeEveryMatchStub).TakeMatch(s.starredMatch))
}

func (s MatchesFilterSuite) TestTeamFilterSkipUknownTeam() {
	s.False(NewTeamFilter(s.takeEveryMatchStub, "team").TakeMatch(s.todayMatch))
}

func (s MatchesFilterSuite) TestTeamFilterTakeMatchWithGivenTeam() {
	s.True(NewTeamFilter(s.takeEveryMatchStub, "first").TakeMatch(s.todayMatch))
	s.True(NewTeamFilter(s.takeEveryMatchStub, "second").TakeMatch(s.todayMatch))
}

func (s MatchesFilterSuite) TestFiltersAlreadySkippedMatches() {
	skip := &RejectingFilterStub{}
	s.False(NewFilterTodayMatches(skip).TakeMatch(s.todayMatch))
	s.False(NewStarredFilter(skip).TakeMatch(s.todayMatch))
	s.False(NewTeamFilter(skip, "team").TakeMatch(s.todayMatch))
}

func TestMatchesFilterSuite(t *testing.T) {
	suite.Run(t, new(MatchesFilterSuite))
}

type RejectingFilterStub struct{}

func (f *RejectingFilterStub) TakeMatch(m Match) bool { return false }
