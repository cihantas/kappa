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
	opts := &twitch.UsersGetOptions{
		Login: []string{"devvv", "killinginpink"},
	}
	users, resp, err := client.Users.Get(opts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Users: %v", users)
	fmt.Printf("Response: %v", resp)
}
