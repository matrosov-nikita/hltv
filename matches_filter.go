package main

import "strings"

type MatchesFilter interface {
	TakeMatch(m Match) bool
}

// FilterTodayMatches takes only today matches.
type FilterTodayMatches struct {
	filter MatchesFilter
}

func NewFilterTodayMatches(filter MatchesFilter) MatchesFilter {
	return &FilterTodayMatches{filter}
}

func (f *FilterTodayMatches) TakeMatch(m Match) bool {
	return f.filter.TakeMatch(m) && m.IsToday()
}

// StarredFilter takes only starred matches.
type StarredFilter struct {
	filter MatchesFilter
}

func NewStarredFilter(filter MatchesFilter) MatchesFilter {
	return &StarredFilter{filter}
}

// TakeMatch takes matches with given team.
func (f *StarredFilter) TakeMatch(m Match) bool {
	return f.filter.TakeMatch(m) && m.Stars > 0
}

type TeamFilter struct {
	filter MatchesFilter
	team   string
}

func NewTeamFilter(filter MatchesFilter, team string) MatchesFilter {
	return &TeamFilter{filter: filter, team: strings.ToLower(team)}
}

func (f *TeamFilter) TakeMatch(m Match) bool {
	return f.filter.TakeMatch(m) &&
		(strings.Contains(strings.ToLower(m.FirstTeam), f.team) ||
			strings.Contains(strings.ToLower(m.SecondTeam), f.team))
}

// TakeEveryMatchFilter default filter which takes all matches.
type TakeEveryMatchFilter struct{}

func (f *TakeEveryMatchFilter) TakeMatch(m Match) bool { return true }
