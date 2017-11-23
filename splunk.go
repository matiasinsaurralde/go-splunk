package splunk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const (
	defaultPath      = "/services/collector/event/1.0"
	authHeaderName   = "authorization"
	authHeaderPrefix = "Splunk "
)

var (
	errInvalidSettings = errors.New("Empty settings")
)

// Client contains Splunk client methods.
type Client struct {
	Token         string
	Endpoint      string
	TLSSkipVerify bool

	httpClient *http.Client
}

// New initializes a new client.
func New(token string, endpoint string, skipVerify bool) (c *Client, err error) {
	if token == "" || endpoint == "" {
		return c, errInvalidSettings
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return c, err
	}
	// Append the default collector API path:
	u.Path = defaultPath
	c = &Client{
		Token:      token,
		Endpoint:   u.String(),
		httpClient: http.DefaultClient,
	}
	if skipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.httpClient = &http.Client{Transport: tr}
	}
	return c, nil
}

// Send sends an event to the Splunk HTTP Event Collector interface.
func (c *Client) Send(event map[string]interface{}) (*http.Response, error) {
	eventWrap := struct {
		Event map[string]interface{} `json:"event"`
	}{Event: event}
	eventJSON, err := json.Marshal(eventWrap)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(eventJSON)
	req, err := http.NewRequest("POST", c.Endpoint, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Add(authHeaderName, authHeaderPrefix+c.Token)
	return c.httpClient.Do(req)
}
