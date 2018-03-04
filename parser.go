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

	// TODO: how to meake this error happen?
	doc, _ := goquery.NewDocumentFromReader(hltvContent)
	p.getMatchesBody(doc).Each(p.parseMatch)
	return p.matches, nil
}

func (p *Parser) getMatchesBody(doc *goquery.Document) *goquery.Selection {
	return doc.Find(".upcoming-matches .match-day > a")
}

func (p *Parser) parseMatch(i int, match *goquery.Selection) {
	m := Match{}
	m.Link, _ = match.Attr("href")

	teams := match.Find(".team")
	m.FirstTeam = teams.Eq(0).Text()
	m.SecondTeam = teams.Eq(1).Text()

	t, _ := match.Find("div.time").Attr("data-unix")
	unixMs, _ := strconv.ParseInt(t, 10, 64)
	m.Time = time.Unix(unixMs/1000, 0)

	m.Stars = match.Find("div.stars").Children().Length()
	p.matches = append(p.matches, m)
}

func (p *Parser) getMatches() []Match {
	return p.matches
}
