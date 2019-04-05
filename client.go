package main

import "net/http"

// ClientWithHeaders represents a HTTP client with custom request headers.
type ClientWithHeaders struct {
	headers map[string]string
}

// NewClientWithHeaders creates a HTTP client which will add given headers to request.
func NewClientWithHeaders(headers map[string]string) *ClientWithHeaders {
	return &ClientWithHeaders{headers: headers}
}

// Get requests given url.
func (c ClientWithHeaders) Get(url string) (*http.Response, error) {
	req, err := c.buildRequest(url)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func (c ClientWithHeaders) buildRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for header, value := range c.headers {
		req.Header.Set(header, value)
	}

	return req, nil
}