package twitch

import "net/url"

// UsersGetCall represents a GET request to the /users endpoint.
type UsersFollowsGetRequest struct {
	client    *Client
	urlParams url.Values
}

// UserFollow represents a follow relationship between to Twitch users.
type UserFollow struct {
	FollowedAt string `json:"followed_at,omitempty"`
	FromID     string `json:"from_id,omitempty"`
	ToID       string `json:"to_id, omitempty"`
}

// UsersFollowsGetResponse represents a response returned by the endpoint /users/follows.
type UsersFollowsGetResponse struct {
	Data       []UserFollow `json:"data"`
	Pagination Pagination
}

// Get returns an instance of UsersFollowsGetCall.
func (s *UsersService) GetFollows() *UserFollow {
	return nil
}
