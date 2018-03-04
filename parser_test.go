package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ParserSuite struct {
	suite.Suite

	client *SpyHTTPClient
	parser *Parser
}

func (s *ParserSuite) SetupTest() {
	f, err := os.Open("./hltv_fixture.html")
	if err != nil {
		panic(err)
	}
	response := &SpyHTTPBody{Reader: f}
	s.client = &SpyHTTPClient{code: http.StatusOK, response: response}
	s.parser = NewParser(s.client)
}

func (s *ParserSuite) AssertMatchFields(m Match, sec int64, link, fTeam, sTeam string, stars int) {
	s.Equal(time.Unix(sec, 0), m.Time)
	s.Equal(link, m.Link)
	s.Equal(fTeam, m.FirstTeam)
	s.Equal(sTeam, m.SecondTeam)
	s.Equal(stars, m.Stars)
}

func (s *ParserSuite) ConfigureClient(code int, err error) {
	s.client.code = code
	s.client.err = err
}

func (s *ParserSuite) TestParserReturnErrorWhenClientFails() {
	s.ConfigureClient(http.StatusInternalServerError, errors.New("error"))
	_, err := s.parser.FetchMatches()
	s.NotNil(err)
}

func (s *ParserSuite) TestQueryBuildedProperly() {
	s.parser.FetchMatches()
	s.Equal("https://www.hltv.org/matches", s.client.url)
}

func (s *ParserSuite) TestParserClosesResponse() {
	s.parser.FetchMatches()
	s.True(s.client.response.closeCalled)
}

func (s *ParserSuite) TestParserHandlerErrorCode() {
	s.ConfigureClient(http.StatusNotFound, nil)
	_, err := s.parser.FetchMatches()
	s.NotNil(err)
}

func (s *ParserSuite) TestParserReturnsAllMatches() {
	tt := []struct {
		sec        int64
		link       string
		firstTeam  string
		secondTeam string
		stars      int
	}{
		{sec: 1520168400, link: "/testmatch1", firstTeam: "eXtatus", secondTeam: "Unity", stars: 0},
		{sec: 1520172000, link: "/testmatch2", firstTeam: "GoodJob", secondTeam: "DreamEaters", stars: 5},
		{sec: 1520168400, link: "/testmatch1", firstTeam: "eXtatus", secondTeam: "Unity", stars: 0},
		{sec: 1520172000, link: "/testmatch2", firstTeam: "GoodJob", secondTeam: "DreamEaters", stars: 5},
	}

	matches, _ := s.parser.FetchMatches()

	s.Len(matches, 4)
	for i, m := range matches {
		s.AssertMatchFields(m, tt[i].sec, tt[i].link,
			tt[i].firstTeam, tt[i].secondTeam, tt[i].stars)
	}
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}

type SpyHTTPClient struct {
	url string

	code     int
	err      error
	response *SpyHTTPBody
}

func (c *SpyHTTPClient) Get(url string) (resp *http.Response, err error) {
	c.url = url

	return &http.Response{
		Body:       c.response,
		StatusCode: c.code,
	}, c.err
}

type SpyHTTPBody struct {
	io.Reader

	closeCalled bool
}

func (b *SpyHTTPBody) Close() error {
	b.closeCalled = true
	return nil
}
