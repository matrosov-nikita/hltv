package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {

	onlyStarred := flag.Bool("starred", false, "Show only starred matches")
	onlyToday := flag.Bool("today", false, "Show only today matches")
	onlyForTeam := flag.String("team", "", "Show only for given team")
	flag.Parse()

	parser := NewParser(http.DefaultClient)
	_, err := parser.FetchMatches()
	if err != nil {
		log.Fatalln(err)
	}

	var filter MatchesFilter = parser

	if *onlyStarred {
		filter = NewStarredFilter(filter)
	}

	if *onlyToday {
		filter = NewFilterTodayMatches(filter)
	}

	team := *onlyForTeam
	if len(team) > 0 {
		filter = NewTeamFilter(filter, team)
	}

	for _, m := range filter.getMatches() {
		fmt.Println(m.String())
	}
}
