package twitch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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

// UsersGetResponse ...
type UsersGetResponse struct {
	Data []User `json:"data"`
}

// Get returns an instance of UsersGetCall.
func (s *UsersService) Get() *UsersGetCall {
	c := &UsersGetCall{service: s.service, urlParams: make(map[string][]string)}
	return c
}

// UsersGetCall represents a GET request to the /users endpoint.
type UsersGetCall struct {
	service   *Service
	urlParams url.Values
}

// ID adds a user IDs to the request.
func (c *UsersGetCall) ID(id string) *UsersGetCall {
	c.urlParams.Add("id", id)
	return c
}

// IDs adds multiple user IDs to the request.
func (c *UsersGetCall) IDs(ids []string) *UsersGetCall {
	joinedIDs := strings.Join(ids, ",")
	c.urlParams.Add("id", joinedIDs)
	return c
}

// Login adds a user login to the request.
func (c *UsersGetCall) Login(login string) *UsersGetCall {
	c.urlParams.Add("login", login)
	return c
}

// Logins adds multiple user logins to the request.
func (c *UsersGetCall) Logins(logins []string) *UsersGetCall {
	joinedLogins := strings.Join(logins, ",")
	c.urlParams.Add("id", joinedLogins)
	return c
}

func (c *UsersGetCall) Do() (*UsersGetResponse, error) {
	<-rateLimiter

	// TODO: If no id and login query parameter is specified, try to use the Bearer token.

	reqURL := basePath + "users?" + c.urlParams.Encode()
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	//req.Header.Add("Accept", "application/vnd.twitchtv."+apiVersion+"+json")
	req.Header.Add("Client-ID", c.service.clientID)
	if c.service.accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.service.accessToken)
	}

	res, err := c.service.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 && res.StatusCode < 500 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(body))

		statusCodeString := strconv.Itoa(res.StatusCode)
		fmt.Println()
		return nil, errors.New("Status code not 200, it is " + statusCodeString)
	} else if res.StatusCode >= 500 {
		return nil, errors.New("Status code not 200, it is " + strconv.Itoa(res.StatusCode))
	}

	ur := &UsersGetResponse{}
	if err := json.NewDecoder(res.Body).Decode(ur); err != nil {
		return nil, err
	}

	return ur, nil
}

//
//
//

// ChannelFollowsResponse ...
//type ChannelFollowsResponse struct {
//	Cursor  string `json:"_cursor,omitempty"`
//	Follows []struct {
//		CreatedAt     time.Time `json:"created_at,omitempty"`
//		Notifications bool      `json:"notifications,omitempty"`
//		User          User      `json:"user,omitempty"`
//	} `json:"follows,omitempty"`
//	Total int `json:"_total,omitempty"`
//}
//
//type ChannelFollowsGetCall struct {
//	service   *Service
//	urlParams url.Values
//	channelID string
//}
//
//// ChannelFollowsResponse ...
//type ChannelFollowsResponse struct {
//	Cursor  string `json:"_cursor,omitempty"`
//	Follows []struct {
//		CreatedAt     time.Time `json:"created_at,omitempty"`
//		Notifications bool      `json:"notifications,omitempty"`
//		User          User      `json:"user,omitempty"`
//	} `json:"follows,omitempty"`
//	Total int `json:"_total,omitempty"`
//}
//
//func (r *UsersService) GetFollows() *ChannelFollowsGetCall {
//	c := &ChannelFollowsGetCall{service: r.service, urlParams: make(map[string][]string)}
//	return c
//}
//
//func (c *ChannelFollowsGetCall) ChannelID(channelID string) *ChannelFollowsGetCall {
//	c.channelID = channelID
//	return c
//}
//
//// Limit defines the maximum number of objects to return. Maximum: 100.
//func (c *ChannelFollowsGetCall) Limit(limit int) *ChannelFollowsGetCall {
//	c.urlParams.Add("limit", strconv.Itoa(limit))
//	return c
//}
//
//// Offset defines the object offset for pagination results.
//func (c *ChannelFollowsGetCall) Offset(offset int) *ChannelFollowsGetCall {
//	c.urlParams.Add("offset", strconv.Itoa(offset))
//	return c
//}
//
//// Cursor tells the server where to start fetching the next set of result, in a
//// multi-page response.
//func (c *ChannelFollowsGetCall) Cursor(cursor string) *ChannelFollowsGetCall {
//	c.urlParams.Add("cursor", cursor)
//	return c
//}
//
//// Direction of sorting. Valid values: asc, desc (newest first).
//func (c *ChannelFollowsGetCall) Direction(direction string) *ChannelFollowsGetCall {
//	c.urlParams.Add("direction", direction)
//	return c
//}
//
//// Do ...
//func (c *ChannelFollowsGetCall) Do() (*ChannelFollowsResponse, error) {
//	<-rateLimiter
//
//	url := basePath + "channels/" + c.channelID + "/follows?" + c.urlParams.Encode()
//
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		return nil, err
//	}
//	req.Header.Add("Accept", "application/vnd.twitchtv."+apiVersion+"+json")
//	req.Header.Add("Client-ID", c.service.accessToken)
//
//	res, err := c.service.client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer res.Body.Close()
//
//	if res.StatusCode != 200 {
//		body, err := ioutil.ReadAll(res.Body)
//		if err != nil {
//			return nil, err
//		}
//		fmt.Println(string(body))
//
//		statusCodeString := strconv.Itoa(res.StatusCode)
//		fmt.Println()
//		return nil, errors.New("Status code not 200, it is " + statusCodeString)
//	}
//
//	cfr := &ChannelFollowsResponse{}
//	if err := json.NewDecoder(res.Body).Decode(cfr); err != nil {
//		return nil, err
//	}
//
//	return cfr, nil
//}
