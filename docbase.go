package docbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.docbase.io/teams/%s"
	apiVersion     = "2"
	userAgent      = "DocBase Go %s"

	// https://help.docbase.io/posts/45703#利用制限
	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-RateLimit-Remaining"
	headerRateReset     = "X-RateLimit-Reset"
)

const (
	publicScope  = "public"
	groupScope   = "group"
	privateScope = "private"
)

type Client struct {
	BaseURL     *url.URL
	AccessToken string
	Team        string
	Client      *http.Client

	Posts       PostService
	Users       UserService
	Groups      GroupService
	GroupUsers  GroupUserService
	Tags        TagService
	Comments    CommentService
	Attachments AttachmentService
}

// Response is http response wrapper for Dobase
type Response struct {
	*http.Response
	Rate
}

func newResponse(r *http.Response) *Response {
	res := &Response{Response: r}
	res.Rate = parseRate(r)
	return res
}

func parseRate(r *http.Response) Rate {
	var (
		rate Rate
		err  error
	)

	if limit := r.Header.Get(headerRateLimit); limit != "" {
		rate.Limit, err = strconv.Atoi(limit)
		if err != nil {
			rate.err = err
		}
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		rate.Remaining, err = strconv.Atoi(remaining)
		if err != nil {
			rate.err = err
		}
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		v, e := strconv.ParseInt(reset, 10, 64)
		if e != nil {
			rate.err = e
		} else if v != 0 {
			//rate.Reset = Timestamp{time.Unix(v, 0)}
		}
	}
	return rate
}

type Rate struct {
	// The number of requests per hour the client is currently limited to.
	Limit int `json:"limit"`

	// The number of remaining requests the client can make this hour.
	Remaining int `json:"remaining"`

	// The time at which the current rate limit will reset.
	Reset time.Time `json:"reset"`

	err error
}

type ErrorResponse struct {
	Messages []string `json:"messages"`
}

func (e *ErrorResponse) Error() string {
	return strings.Join(e.Messages, "\n - ")
}

// NewClient returns client API
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

	cli.Posts = &PostCli{cli}
	cli.Groups = &GroupCli{cli}
	cli.Users = &UserCli{cli}
	cli.Comments = &CommentCli{cli}
	cli.Tags = &TagCli{cli}
	cli.Attachments = &AttachmentCli{cli}
	cli.GroupUsers = &GroupUserCli{cli}

	return cli
}

type Option func(client *Client)

func OptionDocbaseURL(url *url.URL) Option {
	return func(client *Client) {
		client.BaseURL = url
	}
}

// NewRequest creates a API request with HTTP method, endpoint path and payload
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

// Do sends request and returns API response
func (c *Client) Do(r *http.Request, v interface{}) (*Response, error) {
	resp, err := c.Client.Do(r)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	if err != nil {
		return response, err
	}
	return response, nil
}

func (c *Client) DoUpload(r *http.Request) (FileContent, *Response, error) {
	resp, err := c.Client.Do(r)

	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return nil, response, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, response, err
	}
	return body, response, nil
}

// CheckResponse checks response for errors
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
