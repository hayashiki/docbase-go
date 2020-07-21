package docbase

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.docbase.io/teams/%s"
)

type Client struct {
	BaseURL *url.URL
	AccessToken string
	Team string
	Client *http.Client
}

type Service struct {
	client *Client
}

func NewClient(httpClient *http.Client, team, token string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, err := url.Parse(fmt.Sprintf(defaultBaseURL, team))

	if err != nil {
		log.Fatal(err)
	}

	cli := &Client{
		BaseURL: baseURL,
		AccessToken: token,
		Team: team,
		Client: httpClient,
	}

	return cli
}
