package main

import "strings"

type MatchesFilter interface {
	TakeMatch(m Match) bool
}

type FilterTodayMatches struct {
	filter MatchesFilter
}

func NewFilterTodayMatches(filter MatchesFilter) MatchesFilter {
	return &FilterTodayMatches{filter}
}

func (f *FilterTodayMatches) TakeMatch(m Match) bool {
	return f.filter.TakeMatch(m) && m.IsToday()
}

type StarredFilter struct {
	filter MatchesFilter
}

func NewStarredFilter(filter MatchesFilter) MatchesFilter {
	return &StarredFilter{filter}
}

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

type TakeEveryMatchFilter struct{}

func (f *TakeEveryMatchFilter) TakeMatch(m Match) bool { return true }
