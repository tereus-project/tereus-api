package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xanzy/go-gitlab"
)

type GitlabService struct {
	ClientId     string
	clientSecret string
}

func NewGitlabService(clientId string, clientSecret string) (*GitlabService, error) {
	return &GitlabService{
		ClientId:     clientId,
		clientSecret: clientSecret,
	}, nil
}

type GitlabAccessTokenResponseBody struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int64  `json:"created_at"`

	Error string `json:"error_description"`
}

func (s *GitlabService) GenerateAccessTokenFromCode(code string, redirectUri string) (*GitlabAccessTokenResponseBody, error) {
	url := fmt.Sprintf("https://gitlab.com/oauth/token?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=%s", s.ClientId, s.clientSecret, code, redirectUri)

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

	var body GitlabAccessTokenResponseBody
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		return nil, err
	}

	if body.Error != "" {
		return nil, fmt.Errorf("%s", body.Error)
	}

	return &body, nil
}

type GitlabClient struct {
	client *gitlab.Client
}

func (s *GitlabService) NewClient(accessToken string) (*GitlabClient, error) {
	client, err := gitlab.NewOAuthClient(accessToken)
	if err != nil {
		return nil, err
	}

	return &GitlabClient{
		client: client,
	}, nil
}

func (c *GitlabClient) GetUser() (*gitlab.User, error) {
	user, _, err := c.client.Users.CurrentUser()
	return user, err
}
