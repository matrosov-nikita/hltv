package main

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const HLTVMatchPage = "https://www.hltv.org/matches"

var ErrCannotLoadMatchesPage = errors.New("cannot load matches page")

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type Parser struct {
	client HTTPClient
}

func NewParser(client HTTPClient) *Parser {
	return &Parser{client: client}
}

func (p *Parser) FetchMatches() ([]Match, error) {
	matchesPage, err := p.loadHLTVMatchesPage()
	if err != nil {
		return nil, err
	}

	defer matchesPage.Close()

	var matches []Match

	p.getMatchesHTML(matchesPage).Each(func(i int, matchHTML *goquery.Selection) {
		matches = append(matches, newMatchParser(matchHTML).Parse())
	})

	return matches, nil
}

func (p *Parser) loadHLTVMatchesPage() (io.ReadCloser, error) {
	resp, err := p.client.Get(HLTVMatchPage)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, ErrCannotLoadMatchesPage
	}

	return resp.Body, nil
}

func (p *Parser) getMatchesHTML(page io.Reader) *goquery.Selection {
	doc, _ := goquery.NewDocumentFromReader(page)
	return doc.Find(".upcoming-matches .match-day > a")
}

type matchParser struct {
	doc *goquery.Selection
}

func newMatchParser(doc *goquery.Selection) *matchParser {
	return &matchParser{doc: doc}
}

func (p *matchParser) Parse() Match {
	firstTeam, secondTeam := p.Teams()
	return NewMatch(firstTeam, secondTeam, p.Link(), p.Stars(), p.Time())
}

func (p *matchParser) Link() string {
	link, _ := p.doc.Attr("href")
	return link
}

func (p *matchParser) Stars() int {
	return p.doc.Find("div.stars").Children().Length()
}

func (p *matchParser) Teams() (string, string) {
	teams := p.doc.Find(".team")
	return teams.Eq(0).Text(), teams.Eq(1).Text()
}

func (p *matchParser) Time() time.Time {
	t, _ := p.doc.Find("div.time").Attr("data-unix")
	unixMs, _ := strconv.ParseInt(t, 10, 64)
	return time.Unix(unixMs/1000, 0)
}
