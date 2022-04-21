package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/v43/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type GithubService struct {
	ClientId     string
	clientSecret string
}

func NewGithubService(clientId string, clientSecret string) (*GithubService, error) {
	return &GithubService{
		ClientId:     clientId,
		clientSecret: clientSecret,
	}, nil
}

type GithubAccessTokenResponseBody struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (s *GithubService) GenerateAccessTokenFromCode(code string) (*GithubAccessTokenResponseBody, error) {
	url := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", s.ClientId, s.clientSecret, code)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logrus.Debug(string(rawBody))

	var body GithubAccessTokenResponseBody
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

type GithubClient struct {
	client *github.Client
}

func (s *GithubService) NewClient(accessToken string) *GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return &GithubClient{
		client: github.NewClient(tc),
	}
}

func (c *GithubClient) GetUser() (*github.User, error) {
	user, _, err := c.client.Users.Get(context.Background(), "")
	return user, err
}
