package main

import (
	"net/http"
	"time"

	"os"

	"fmt"

	"github.com/cihantas/kappa/twitch"
)

func main() {
	at := &twitch.AuthenticatedTransport{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
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
	_, resp, err := client.Users.Get(opts)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Users: %v", users)
	fmt.Printf("RateLimit-Remaining: %s\n", resp.RawResponse.Header.Get("RateLimit-Remaining"))
	fmt.Printf("Status Code: %d\n", resp.RawResponse.StatusCode)
}
