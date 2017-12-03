package twitch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

	r := &UsersFollowsGetResponse{}
	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return nil, nil, err
	}

	return &r.Data, r, nil
	// TODO: If no id and login query parameter is specified, try to use the Bearer token.
}
