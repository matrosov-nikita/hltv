package main

import "strings"

type MatchesFilter interface {
	getMatches() []Match
}

type FilterTodayMatches struct {
	filter MatchesFilter
}

func NewFilterTodayMatches(filter MatchesFilter) MatchesFilter {
	return &FilterTodayMatches{filter}
}

func (f *FilterTodayMatches) getMatches() []Match {
	var todayMatches []Match

	for _, match := range f.filter.getMatches() {
		if match.IsToday() {
			todayMatches = append(todayMatches, match)
		}
	}

	return todayMatches
}

type StarredFilter struct {
	filter MatchesFilter
}

func NewStarredFilter(filter MatchesFilter) MatchesFilter {
	return &StarredFilter{filter}
}

func (f *StarredFilter) getMatches() []Match {
	var starredMatches []Match

	for _, match := range f.filter.getMatches() {
		if match.Stars > 0 {
			starredMatches = append(starredMatches, match)
		}
	}

	return starredMatches
}

type TeamFilter struct {
	filter MatchesFilter
	team   string
}

func NewTeamFilter(filter MatchesFilter, team string) MatchesFilter {
	return &TeamFilter{filter: filter, team: strings.ToLower(team)}
}

func (f *TeamFilter) getMatches() []Match {
	var matchesForTeam []Match

	for _, match := range f.filter.getMatches() {
		if strings.Contains(strings.ToLower(match.FirstTeam), f.team) ||
			strings.Contains(strings.ToLower(match.SecondTeam), f.team) {
			matchesForTeam = append(matchesForTeam, match)
		}
	}

	return matchesForTeam
}
