package twitch

import (
	"fmt"
	"net/http"
	"testing"
)

var client = &http.Client{}
var clientID = ""

func TestUsersService_Get(t *testing.T) {
	service, err := New(client, clientID)
	if err != nil {
		t.Error(err)
	}

	call := service.Users.Get()
	if call == nil {
		t.Failed()
	}
}

func TestChannelFollowsGetCall_Do(t *testing.T) {
	service, err := New(client, clientID)
	if err != nil {
		t.Error(err)
	}

	call := service.Users.Get()
	resp, err := call.ID("44322889").Do()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", resp)
}
