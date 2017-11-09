package twitch

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	basePath = "https://api.twitch.tv/helix/"

	headerRateLimit          = "RateLimit-Limit"
	headerRateLimitRemaining = "RateLimit-Remaining"
	headerRateLimitReset     = "RateLimit-Reset"
)

var (
	errTooManyRequests = errors.New("Rate-limit exceeded")
)

type Error struct {
	StatusHuman string `json:"error"`
	Status      int    `json:"status"`
	Message     string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%v %v error caused by %v",
		e.Status, e.StatusHuman, e.Message)
}

// This rate limiter (1) should be used at the beginning of each method making a request to
// the Twitch API to ensure a pause of 500ms between each request.
// According to the Twitch API documentation we are allowed to make 2 requests per second (2).
//
// (1) Rate-Limiting Example: https://gobyexample.com/rate-limiting
// (2) Twitch API Documentation: https://dev.twitch.tv/docs/api#rate-limit
var rateLimiter = time.Tick(500 * time.Millisecond)

// New returns a new service instance.
// If no httpClient is provided, http.DefaultClient will be used.
func New(client *http.Client, clientID string) (*Service, error) {
	if client == nil {
		client = http.DefaultClient
	}

	s := &Service{client: client, basePath: basePath, clientID: clientID}
	s.Users = NewUsersService(s)

	return s, nil
}

// Service represents the Twitch API.
type Service struct {
	client      *http.Client
	basePath    string
	clientID    string
	accessToken string

	Users *UsersService
}

func (s *Service) AccessToken(accessToken string) *Service {
	s.accessToken = accessToken
	return s
}

func NewUsersService(s *Service) *UsersService {
	us := &UsersService{service: s}
	return us
}

// var client = &http.Client{Timeout: 20 * time.Second}
