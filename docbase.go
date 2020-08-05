package docbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL = "https://api.docbase.io/teams/%s"
	apiVersion     = "2"
	userAgent = "DocBase Go %s"
)

const (
	publicScope = "public"
	groupScope = "group"
	privateScope = "private"
)

type Client struct {
	BaseURL     *url.URL
	AccessToken string
	Team        string
	Client      *http.Client

	Posts    *PostService
	Users    *UserService
	Groups   *GroupService
	Tags     *TagService
	Comments *CommentService
	Attachments *AttachmentService
}

type ErrorResponse struct {
	Messages []string `json:"messages"`
}

func (e *ErrorResponse) Error() string {
	return strings.Join(e.Messages, "\n - ")
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

	cli.Posts = NewPostService(cli)
	cli.Groups = NewGroupService(cli)
	cli.Users = NewUserService(cli)
	cli.Comments = NewCommentService(cli)
	cli.Tags = NewTagService(cli)
	cli.Attachments = NewAttachmentService(cli)

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

	//hoge := strings.NewReader(buf)
	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(buf))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-DocBaseToken", c.AccessToken)
	req.Header.Add("X-Api-Version", apiVersion)
	req.Header.Add("USER_AGENT", userAgent)

	return req, nil
}

func (c *Client) Do(r *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.Client.Do(r)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (c *Client) DoBinary(r *http.Request) (FileContent, *http.Response, error) {
	resp, err := c.Client.Do(r)

	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return nil, resp, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, resp, err
	}
	return body, resp, nil
}

func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusCreated:
		return nil
	case http.StatusInternalServerError:
		return &ErrorResponse{
			Messages: []string{"Internal Server Error"},
		}
	case http.StatusBadRequest:
		return &ErrorResponse{
			Messages: []string{"Bad Request"},
		}
	case http.StatusForbidden:
		return &ErrorResponse{
			Messages: []string{"Forbidden"},
		}
	default:
		var errResp ErrorResponse
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&errResp)
		if err != nil {
			errResp.Messages = []string{"Couldn't decode response body JSON"}
		}
		return &errResp
	}
}
