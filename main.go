package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	onlyStarred := flag.Bool("starred", false, "Show only starred matches")
	onlyToday := flag.Bool("today", false, "Show only today matches")
	onlyForTeam := flag.String("team", "", "Show only for given team")
	flag.Parse()

	parser := NewParser(NewClientWithHeaders(map[string]string{
		"User-Agent": "bot",
	}))

	matches, err := parser.FetchMatches()
	if err != nil {
		log.Fatalln(err)
	}

	var filter MatchesFilter = &TakeEveryMatchFilter{}

	if *onlyStarred {
		filter = NewStarredFilter(filter)
	}

	if *onlyToday {
		filter = NewFilterTodayMatches(filter)
	}

	if len(*onlyForTeam) > 0 {
		filter = NewTeamFilter(filter, *onlyForTeam)
	}

	for _, m := range matches {
		if filter.TakeMatch(m) {
			fmt.Println(m.String())
		}
	}
}
