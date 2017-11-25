package twitch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"net/url"

	"github.com/google/go-querystring/query"
)

// UsersService handles communication with the user related methods of the Twitch API.
type UsersService struct {
	client *Client
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

// ...
type UsersGetOptions struct {
	ID    []string `url:"id,omitempty"`
	Login []string `url:"login,omitempty"`
}

// UsersGetResponse ...
type UsersGetResponse struct {
	Data []User `json:"data"`
}

// Get returns an instance of UsersGetCall.
func (s *UsersService) Get(opt *UsersGetOptions) (*[]User, *http.Response, error) {
	// Build url.
	q, err := query.Values(opt)
	if err != nil {
		return nil, nil, err
	}
	u, err := url.Parse("users")
	if err != nil {
		return nil, nil, err
	}
	u.RawQuery = q.Encode()

	req, err := s.client.NewRequest("GET", u.String())
	if err != nil {
		return nil, nil, err
	}

	res, err := s.client.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	// TODO
	if res.StatusCode != 200 && res.StatusCode < 500 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, nil, err
		}
		fmt.Println(string(body))

		statusCodeString := strconv.Itoa(res.StatusCode)
		fmt.Println()
		return nil, nil, errors.New("Status code not 200, it is " + statusCodeString)
	} else if res.StatusCode >= 500 {
		return nil, nil, errors.New("Status code not 200, it is " + strconv.Itoa(res.StatusCode))
	}

	ur := &UsersGetResponse{}
	if err := json.NewDecoder(res.Body).Decode(ur); err != nil {
		return nil, nil, err
	}

	return &ur.Data, res, nil
	// TODO: If no id and login query parameter is specified, try to use the Bearer token.
}
