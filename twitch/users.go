package twitch

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

// ...
type UsersGetResponse struct {
	Data        []User `json:"data"`
	RawResponse *http.Response
}

// Get returns an instance of UsersGetCall.
func (s *UsersService) Get(opt *UsersGetOptions) (*[]User, *UsersGetResponse, error) {
	u, err := buildURL("users", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest("GET", u.String())
	if err != nil {
		return nil, nil, err
	}

	res, err := s.client.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if err = checkResponse(res); err != nil {
		return nil, nil, err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	ugr := &UsersGetResponse{}
	if err := json.Unmarshal(bodyBytes, ugr); err != nil {
		return nil, nil, err
	}
	ugr.RawResponse = res

	return &ugr.Data, ugr, nil
	// TODO: If no id and login query parameter is specified, try to use the Bearer token.
}
