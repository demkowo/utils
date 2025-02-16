package httpclient

import (
	"bytes"
	"io"
	"net/http"
)

var (
	client HTTPClient
)

type HTTPClient interface {
	Get(url string, headers map[string]string) (*http.Response, error)
	Post(url string, body []byte, headers map[string]string) (*http.Response, error)
	Put(url string, body []byte, headers map[string]string) (*http.Response, error)
	Patch(url string, body []byte, headers map[string]string) (*http.Response, error)
	Delete(url string, headers map[string]string) (*http.Response, error)
	Head(url string, headers map[string]string) (*http.Response, error)
	Options(url string, headers map[string]string) (*http.Response, error)
}

func NewClient() HTTPClient {
	if isMock {
		client = &clientMock{mocks: make(map[string]Mock)}
	} else {
		client = &cli{httpClient: &http.Client{}}
	}
	return client
}

type cli struct {
	httpClient *http.Client
}

func (c *cli) request(method, url string, body []byte, headers map[string]string) (*http.Response, error) {
	var rdr io.Reader
	if len(body) > 0 {
		rdr = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, rdr)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.httpClient.Do(req)
}

func (c *cli) Get(url string, headers map[string]string) (*http.Response, error) {
	return c.request(http.MethodGet, url, nil, headers)
}

func (c *cli) Post(url string, body []byte, headers map[string]string) (*http.Response, error) {
	return c.request(http.MethodPost, url, body, headers)
}

func (c *cli) Put(url string, body []byte, headers map[string]string) (*http.Response, error) {
	return c.request(http.MethodPut, url, body, headers)
}

func (c *cli) Patch(url string, body []byte, headers map[string]string) (*http.Response, error) {
	return c.request(http.MethodPatch, url, body, headers)
}

func (c *cli) Delete(url string, headers map[string]string) (*http.Response, error) {
	return c.request(http.MethodDelete, url, nil, headers)
}

func (c *cli) Head(url string, headers map[string]string) (*http.Response, error) {
	return c.request(http.MethodHead, url, nil, headers)
}

func (c *cli) Options(url string, headers map[string]string) (*http.Response, error) {
	return c.request(http.MethodOptions, url, nil, headers)
}
