package twitch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	basePath   = "https://api.twitch.tv/helix/"
	apiVersion = "v5"

	headerRateLimit          = "RateLimit-Limit"
	headerRateLimitRemaining = "RateLimit-Remaining"
	headerRateLimitReset     = "RateLimit-Reset"
)

var (
	errTooManyRequests = errors.New("Rate-limit exceeded")
)

// This rate limiter (1) should be used at the beginning of each method making a request to
// the Twitch API to ensure a pause of 500ms between each request.
// According to the Twitch API documentation we are allowed to make 2 requests per second (2).
//
// Source (1): https://gobyexample.com/rate-limiting
// Source (2): https://dev.twitch.tv/docs/api#rate-limit
var rateLimiter = time.Tick(500 * time.Millisecond)

// New returns a new service instance.
func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	s := &Service{client: client, basePath: basePath}
	s.Users = NewUsersService(s)

	return s, nil
}

// Service represents the Twitch API.
type Service struct {
	client      *http.Client
	basePath    string
	clientID    string
	accessToken string

	Users *UsersService
}

func (s *Service) AccessToken(accessToken string) *Service {
	s.accessToken = accessToken
	return s
}

func NewUsersService(s *Service) *UsersService {
	us := &UsersService{service: s}
	return us
}

// UsersService handles communication with the user related methods of the Twitch API.
type UsersService struct {
	service *Service
}

// User represents a Twitch user.
type User struct {
	BroadcasterType string `json:"broadcaster_type,omitempty"`
	Description     string `json:"description,omitempty"`
	DisplayName     string `json:"display_name,omitempty"`
	Email           string `json:"email,omitempty"`
	ID              string `json:"id,omitempty"`
	Login           string `json:"login,omitempty"`
	OfflineImageURL string `json:"offline_image_url,omitempty"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	Type            string `json:"video_banner,omitempty"`
	ViewCount       int    `json:"view_count,omitempty"`
}

type ChannelFollowsGetCall struct {
	service   *Service
	urlParams url.Values
	channelID string
}

// ChannelFollowsResponse ...
type ChannelFollowsResponse struct {
	Cursor  string `json:"_cursor,omitempty"`
	Follows []struct {
		CreatedAt     time.Time `json:"created_at,omitempty"`
		Notifications bool      `json:"notifications,omitempty"`
		User          User      `json:"user,omitempty"`
	} `json:"follows,omitempty"`
	Total int `json:"_total,omitempty"`
}

func (r *UsersService) Get() *ChannelFollowsGetCall {
	c := &ChannelFollowsGetCall{service: r.service, urlParams: make(map[string][]string)}
	return c
}

func (c *ChannelFollowsGetCall) ChannelID(channelID string) *ChannelFollowsGetCall {
	c.channelID = channelID
	return c
}

// Limit defines the maximum number of objects to return. Maximum: 100.
func (c *ChannelFollowsGetCall) Limit(limit int) *ChannelFollowsGetCall {
	c.urlParams.Add("limit", strconv.Itoa(limit))
	return c
}

// Offset defines the object offset for pagination results.
func (c *ChannelFollowsGetCall) Offset(offset int) *ChannelFollowsGetCall {
	c.urlParams.Add("offset", strconv.Itoa(offset))
	return c
}

// Cursor tells the server where to start fetching the next set of result, in a
// multi-page response.
func (c *ChannelFollowsGetCall) Cursor(cursor string) *ChannelFollowsGetCall {
	c.urlParams.Add("cursor", cursor)
	return c
}

// Direction of sorting. Valid values: asc, desc (newest first).
func (c *ChannelFollowsGetCall) Direction(direction string) *ChannelFollowsGetCall {
	c.urlParams.Add("direction", direction)
	return c
}

// Do ...
func (c *ChannelFollowsGetCall) Do() (*ChannelFollowsResponse, error) {
	<-rateLimiter

	url := basePath + "channels/" + c.channelID + "/follows?" + c.urlParams.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.twitchtv."+apiVersion+"+json")
	req.Header.Add("Client-ID", c.service.accessToken)

	res, err := c.service.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(body))

		statusCodeString := strconv.Itoa(res.StatusCode)
		fmt.Println()
		return nil, errors.New("Status code not 200, it is " + statusCodeString)
	}

	cfr := &ChannelFollowsResponse{}
	if err := json.NewDecoder(res.Body).Decode(cfr); err != nil {
		return nil, err
	}

	return cfr, nil
}

// var client = &http.Client{Timeout: 20 * time.Second}
