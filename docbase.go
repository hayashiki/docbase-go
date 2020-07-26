package docbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.docbase.io/teams/%s"
	apiVersion     = "2"
)

type Client struct {
	BaseURL     *url.URL
	AccessToken string
	Team        string
	Client      *http.Client
}

type Service struct {
	client *Client
}

func NewClient(httpClient *http.Client, team, token string, opts ...Option) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, err := url.Parse(fmt.Sprintf(defaultBaseURL, team))

	if err != nil {
		log.Fatal(err)
	}

	cli := &Client{
		AccessToken: token,
		Team:        team,
		Client:      httpClient,
	}

	for _, opt := range opts {
		opt(cli)
	}

	if cli.BaseURL == nil {
		cli.BaseURL = baseURL
	}

	return cli
}

type Option func(client *Client)

func OptionDocbaseURL(url *url.URL) Option {
	return func(client *Client) {
		client.BaseURL = url
	}
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {

	u, err := url.Parse(fmt.Sprintf("%s/%s", c.BaseURL.String(), path))

	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(buf))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-DocBaseToken", c.AccessToken)
	req.Header.Add("X-Api-Version", apiVersion)

	return req, nil
}

func (c *Client) Do(r *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.Client.Do(r)


	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&v)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
