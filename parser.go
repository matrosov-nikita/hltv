package main

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var ErrCannotGetHltv = errors.New("hltv returned error")

type Parser struct {
	client HTTPClient

	matches []Match
}

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

func NewParser(client HTTPClient) *Parser {
	return &Parser{client: client}
}

func (p *Parser) loadHltvPage() (io.ReadCloser, error) {
	resp, err := p.client.Get("https://www.hltv.org/matches")
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, ErrCannotGetHltv
	}

	return resp.Body, nil
}

func (p *Parser) FetchMatches() ([]Match, error) {
	hltvContent, err := p.loadHltvPage()
	if err != nil {
		return nil, err
	}

	defer hltvContent.Close()

	doc, _ := goquery.NewDocumentFromReader(hltvContent)

	p.getMatchesDocuments(doc).Each(func(i int, doc *goquery.Selection) {
		p.matches = append(p.matches, newMatchParser(doc).Parse())
	})

	return p.matches, nil
}

func (p *Parser) getMatchesDocuments(doc *goquery.Document) *goquery.Selection {
	return doc.Find(".upcoming-matches .match-day > a")
}

type matchParser struct {
	doc *goquery.Selection
}

func newMatchParser(doc *goquery.Selection) *matchParser {
	return &matchParser{doc: doc}
}

func (p *matchParser) Parse() Match {
	link, _ := p.doc.Attr("href")

	firstTeam, secondTeam := p.parseTeams()
	unixTime := p.parseTime()
	stars := p.doc.Find("div.stars").Children().Length()

	return NewMatch(firstTeam, secondTeam, link, stars, unixTime)
}

func (p *matchParser) parseTeams() (string, string) {
	teams := p.doc.Find(".team")
	return teams.Eq(0).Text(), teams.Eq(1).Text()
}

func (p *matchParser) parseTime() time.Time {
	t, _ := p.doc.Find("div.time").Attr("data-unix")
	unixMs, _ := strconv.ParseInt(t, 10, 64)
	return time.Unix(unixMs/1000, 0)
}
