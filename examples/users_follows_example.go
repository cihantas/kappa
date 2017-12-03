package main

import (
	"net/http"
	"time"

	"fmt"

	"github.com/cihantas/kappa/twitch"
)

func main() {
	at := &twitch.AuthenticatedTransport{
		ClientID:     "",
		ClientSecret: "",
		Transport: &http.Transport{
			IdleConnTimeout: 30 * time.Second,
		},
	}
	hc := &http.Client{Transport: at}
	client, err := twitch.NewClient(hc)
	if err != nil {
		panic(err)
	}

	opts := &twitch.UserFollowsGetOptions{}
	follows, resp, err := client.Users.ListFollows(opts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User Follows: %v", follows)
	fmt.Printf("Response: %v", resp)
}
