package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ClientSuite struct {
	suite.Suite
	c *ClientWithHeaders

	testURL string
}

func (s *ClientSuite) SetupTest() {
	s.c = NewClientWithHeaders(map[string]string{
		"TestHeader": "test",
	})
	s.testURL = "http://some_url"
}

func (s *ClientSuite) TestBuildRequest() {
	req, err := s.c.buildRequest(s.testURL)
	s.Nil(err)
	s.Equal(req.Header.Get("TestHeader"), "test")
	s.Equal(req.URL.String(), s.testURL)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}