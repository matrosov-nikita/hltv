package main

import (
	"errors"
	"io"
	"net/http"
	"strings"
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
	response := &SpyHTTPBody{Reader: strings.NewReader(rawHTMLMatches)}
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
	s.Equal(HLTVMatchPage, s.client.url)
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

var rawHTMLMatches = `
<!DOCTYPE html>
<html lang="en">
  <head>
  </head>
  <body>
    <div class="upcoming-matches">
      <div class="" data-zonedgrouping-headline-format="YYYY-MM-dd" data-zonedgrouping-headline-classes="standard-headline" data-zonedgrouping-group-classes="match-day">
          <div class="match-day">
              <span class="standard-headline">2018-03-04</span>
              <a href="/testmatch1" class="a-reset block upcoming-match standard-box" data-zonedgrouping-entry-unix="1520168400000">
                <table class="table">
                  <tr>
                    <td class="time">
                      <div class="time" data-time-format="HH:mm" data-unix="1520168400000">14:00</div>
                    </td>
                    <td class="team-cell">
                      <div class="line-align"><img alt="eXtatus" src="https://static.hltv.org/images/team/logo/6745" class="logo" title="eXtatus">
                        <div class="team">eXtatus</div>
                      </div>
                    </td>
                    <td class="vs">vs</td>
                    <td class="team-cell">
                      <div class="line-align"><img alt="Unity" src="https://static.hltv.org/images/team/logo/7058" class="logo" title="Unity">
                        <div class="team">Unity</div>
                      </div>
                    </td>
                    <td class="event"><img alt="X-BET.co Invitational 2" src="https://static.hltv.org/images/eventLogos/3585.png" class="event-logo" title="X-BET.co Invitational 2"><span class="event-name">X-BET.co Invitational 2</span></td>
                    <td class="star-cell">
                      <div class="map-text">bo3</div>
                    </td>
                  </tr>
                </table>
              </a><a href="/testmatch2" class="a-reset block upcoming-match standard-box" data-zonedgrouping-entry-unix="1520172000000">
                <table class="table">
                  <tr>
                    <td class="time">
                      <div class="time" data-time-format="HH:mm" data-unix="1520172000000">15:00</div>
                    </td>
                    <td class="team-cell">
                      <div class="line-align"><img alt="GoodJob" src="https://static.hltv.org/images/team/logo/7820" class="logo" title="GoodJob">
                        <div class="team">GoodJob</div>
                      </div>
                    </td>
                    <td class="vs">vs</td>
                    <td class="team-cell">
                      <div class="line-align"><img alt="DreamEaters" src="https://static.hltv.org/images/team/logo/8305" class="logo" title="DreamEaters">
                        <div class="team">DreamEaters</div>
                      </div>
                    </td>
                    <td class="event"><img alt="M.Game League " src="https://static.hltv.org/images/eventLogos/3581.png" class="event-logo" title="M.Game League "><span class="event-name">M.Game League </span></td>
                    <td class="star-cell">
                        <div class="map-and-stars">
                            <div class="stars"><i class="fa fa-star star"></i><i class="fa fa-star star"></i><i class="fa fa-star star"></i><i class="fa fa-star star"></i><i class="fa fa-star star"></i></div>
                            <div class="map map-text">bo5</div>
                        </div>
                    </td>
                  </tr>
                </table>
              </a></div>
          <div class="match-day">
                  <span class="standard-headline">2018-03-05</span>
                  <a href="/testmatch1" class="a-reset block upcoming-match standard-box" data-zonedgrouping-entry-unix="1520168400000">
                    <table class="table">
                      <tr>
                        <td class="time">
                          <div class="time" data-time-format="HH:mm" data-unix="1520168400000">14:00</div>
                        </td>
                        <td class="team-cell">
                          <div class="line-align"><img alt="eXtatus" src="https://static.hltv.org/images/team/logo/6745" class="logo" title="eXtatus">
                            <div class="team">eXtatus</div>
                          </div>
                        </td>
                        <td class="vs">vs</td>
                        <td class="team-cell">
                          <div class="line-align"><img alt="Unity" src="https://static.hltv.org/images/team/logo/7058" class="logo" title="Unity">
                            <div class="team">Unity</div>
                          </div>
                        </td>
                        <td class="event"><img alt="X-BET.co Invitational 2" src="https://static.hltv.org/images/eventLogos/3585.png" class="event-logo" title="X-BET.co Invitational 2"><span class="event-name">X-BET.co Invitational 2</span></td>
                        <td class="star-cell">
                          <div class="map-text">bo3</div>
                        </td>
                      </tr>
                    </table>
                  </a><a href="/testmatch2" class="a-reset block upcoming-match standard-box" data-zonedgrouping-entry-unix="1520172000000">
                    <table class="table">
                      <tr>
                        <td class="time">
                          <div class="time" data-time-format="HH:mm" data-unix="1520172000000">15:00</div>
                        </td>
                        <td class="team-cell">
                          <div class="line-align"><img alt="GoodJob" src="https://static.hltv.org/images/team/logo/7820" class="logo" title="GoodJob">
                            <div class="team">GoodJob</div>
                          </div>
                        </td>
                        <td class="vs">vs</td>
                        <td class="team-cell">
                          <div class="line-align"><img alt="DreamEaters" src="https://static.hltv.org/images/team/logo/8305" class="logo" title="DreamEaters">
                            <div class="team">DreamEaters</div>
                          </div>
                        </td>
                        <td class="event"><img alt="M.Game League " src="https://static.hltv.org/images/eventLogos/3581.png" class="event-logo" title="M.Game League "><span class="event-name">M.Game League </span></td>
                        <td class="star-cell">
                            <div class="map-and-stars">
                                <div class="stars"><i class="fa fa-star star"></i><i class="fa fa-star star"></i><i class="fa fa-star star"></i><i class="fa fa-star star"></i><i class="fa fa-star star"></i></div>
                                <div class="map map-text">bo5</div>
                            </div>
                        </td>
                      </tr>
                    </table>
                  </a></div>
      </div>
    </div>
</body>
</html>
`
