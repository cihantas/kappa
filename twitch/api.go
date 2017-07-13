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
	basePath   = "https://api.twitch.tv/kraken/"
	apiVersion = "v5"
)

// This rate limiter should used at the beginning of each method making a request to
// the Twitch API to ensure a pause of 1s between each request.
//
// Source: https://gobyexample.com/rate-limiting
var rateLimiter = time.Tick(1 * time.Second)

// New returns a new service instance.
func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	s := &Service{client: client, basePath: basePath}
	s.Channels = NewChannelsService(s)

	return s, nil
}

type Service struct {
	client      *http.Client
	basePath    string
	clientID    string
	accessToken string

	Channels *ChannelsService
}

func (s *Service) AccessToken(accessToken string) *Service {
	s.accessToken = accessToken
	return s
}

func NewChannelsService(s *Service) *ChannelsService {
	rs := &ChannelsService{service: s}
	return rs
}

type ChannelsService struct {
	service *Service
}

// Channel resource contains information about a Twitch channel.
type ChannelResponse struct {
	BroadcasterLanguage          string    `json:"broadcaster_language,omitempty"`
	BroadcasterType              string    `json:"broadcaster_type,omitempty"`
	CreatedAt                    time.Time `json:"created_at,omitempty"`
	DisplayName                  string    `json:"display_name,omitempty"`
	Email                        string    `json:"email,omitempty"`
	Followers                    int       `json:"followers,omitempty"`
	Game                         string    `json:"game,omitempty"`
	ID                           string    `json:"_id,omitempty"`
	Language                     string    `json:"language,omitempty"`
	Logo                         string    `json:"logo,omitempty"`
	Mature                       bool      `json:"mature,omitempty"`
	Name                         string    `json:"name,omitempty"`
	Partner                      bool      `json:"partner,omitempty"`
	ProfileBanner                string    `json:"profile_banner,omitempty"`
	ProfileBannerBackgroundColor string    `json:"profile_banner_background_color,omitempty"`
	Status                       string    `json:"status,omitempty"`
	StreamKey                    string    `json:"stream_key,omitempty"`
	UpdatedAt                    time.Time `json:"updated_at,omitempty"`
	URL                          string    `json:"url,omitempty"`
	VideoBanner                  string    `json:"video_banner,omitempty"`
	Views                        int       `json:"views,omitempty"`
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

func (r *ChannelsService) Get() *ChannelFollowsGetCall {
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

type User struct {
	Bio         string    `json:"bio,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
	ID          string    `json:"_id,omitempty"`
	Logo        string    `json:"logo,omitempty"`
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// var client = &http.Client{Timeout: 20 * time.Second}
