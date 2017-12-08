package twitch

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// UserFollow represents a follow relationship between to Twitch users.
type UserFollow struct {
	FollowedAt string `json:"followed_at,omitempty"`
	FromID     string `json:"from_id,omitempty"`
	ToID       string `json:"to_id, omitempty"`
}

// ...
type UserFollowsGetOptions struct {
	After  string `url:"after,omitempty"`
	Before string `url:"before,omitempty"`
	First  int    `url:"first,omitempty"`
	FromID string `url:"from_id,omitempty"`
	ToID   string `url:"to_id,omitempty"`
}

// UsersFollowsGetResponse represents a response returned by the endpoint /users/follows.
type UsersFollowsGetResponse struct {
	Data        []UserFollow `json:"data"`
	Pagination  Pagination   `json:"pagination"`
	RawResponse *http.Response
}

// Get returns an instance of UsersGetCall.
func (s *UsersService) ListFollows(opt *UserFollowsGetOptions) (*[]UserFollow, *UsersFollowsGetResponse, error) {
	u, err := buildURL("users/follows", opt)
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

	ufgr := &UsersFollowsGetResponse{}
	if err := json.Unmarshal(bodyBytes, ufgr); err != nil {
		return nil, nil, err
	}
	ufgr.RawResponse = res

	return &ufgr.Data, ufgr, nil
	// TODO: If no id and login query parameter is specified, try to use the Bearer token.
}
