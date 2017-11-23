package twitch

import (
	"net/http"
	"testing"
	"time"
)

func TestNewAuthenticatedClient(t *testing.T) {
	at := &AuthenticatedTransport{
		ClientID:     "test_id",
		ClientSecret: "test_secret",
		Transport: &http.Transport{
			IdleConnTimeout: 30 * time.Second,
		},
	}
	httpClient := &http.Client{Transport: at}
	_, err := NewClient(httpClient)
	if err != nil {
		t.Failed()
	}
}
