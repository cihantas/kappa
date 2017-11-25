package twitch

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL   = "https://api.twitch.tv/helix/"
	defaultUserAgent = "Kappa Helix"

	acceptHeaderHelix = ""
)

// This rate limiter (1) should be used at the beginning of each method making a request to
// the Twitch API to ensure a pause of 500ms between each request.
// According to the Twitch API documentation we are allowed to make 2 requests per second (2).
//
// (1) Rate-Limiting Example: https://gobyexample.com/rate-limiting
// (2) Twitch API Documentation: https://dev.twitch.tv/docs/api#rate-limit
var rateLimiter = time.Tick(500 * time.Millisecond)

// A Client manages communication with the Twitch API.
type Client struct {
	client  *http.Client
	BaseURL *url.URL

	Users *UsersService
}

// New returns a new Twitch API client.
// If no httpClient is provided, http.DefaultClient will be used.
func NewClient(client *http.Client) (*Client, error) {
	if client == nil {
		client = http.DefaultClient
	}

	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{client: client, BaseURL: baseURL}
	c.Users = &UsersService{client: c}

	return c, nil
}

// ...
type AuthenticatedTransport struct {
	ClientID     string
	ClientSecret string
	Transport    http.RoundTripper
}

func (t *AuthenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.ClientID == "" {
		return nil, errors.New("t.ClientID is empty")
	}
	if t.ClientSecret == "" {
		return nil, errors.New("t.ClientSecret is empty")
	}

	// TODO: Copy request as required by RoundTripper interface.
	req.Header.Set("Client-ID", t.ClientID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.ClientID))

	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}

	<-rateLimiter

	return t.Transport.RoundTrip(req)
}

type Pagination struct {
	Cursor string `json:"cursor,omitempty"`
}

func (c *Client) NewRequest(method, urlStr string) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", acceptHeaderHelix)
	//req.Header.Add("Accept", "application/vnd.twitchtv."+apiVersion+"+json")
	req.Header.Set("User-Agent", defaultUserAgent)

	return req, nil
}
