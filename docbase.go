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
	"sync"
	"time"
)

const (
	defaultBaseURL = "https://api.docbase.io/teams/%s"
	apiVersion     = "2"
	userAgent      = "DocBase Go" + version

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

	rateMu    sync.Mutex
	rateLimit Rate

	Posts       PostService
	Users       UserService
	Groups      GroupService
	GroupUsers  GroupUserService
	Tags        TagService
	Comments    CommentService
	Attachments AttachmentService
}

// Response is http response wrapper for DocBase
type Response struct {
	*http.Response
	Rate
	Meta
}

func newResponse(r *http.Response) *Response {
	res := &Response{Response: r}
	res.Rate = parseRate(r)
	return res
}

// NewClient returns client API
func NewClient(httpClient *http.Client, team, token string) *Client {
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
		BaseURL:     baseURL,
	}

	cli.Posts = &postService{cli}
	cli.Groups = &groupService{cli}
	cli.Users = &userService{cli}
	cli.Comments = &commentService{cli}
	cli.Tags = &tagService{cli}
	cli.Attachments = &attachmentService{cli}
	cli.GroupUsers = &groupUserService{cli}

	return cli
}

// NewRequest creates a API request with HTTP method, endpoint path and payload
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {

	u, err := url.Parse(fmt.Sprintf("%s%s", c.BaseURL.String(), path))

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

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-DocBaseToken", c.AccessToken)
	req.Header.Add("X-Api-Version", apiVersion)
	req.Header.Add("USER_AGENT", userAgent)

	return req, nil
}

// Do sends request and returns API response
func (c *Client) Do(r *http.Request, v interface{}) (*Response, error) {

	if err := c.checkRateLimitBeforeDo(r); err != nil {
		return &Response{
			Response: err.Response,
			Rate:     err.Rate,
		}, err
	}

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

	if v == nil {
		return response, nil
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
	if c := r.StatusCode; c == http.StatusOK || c == http.StatusCreated || c == http.StatusNoContent {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		err = json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.err = fmt.Errorf("not json structure")
		}
	}
	switch r.StatusCode {
	case http.StatusTooManyRequests:
		return &RateLimitError{
			Rate:     parseRate(r),
			Response: errorResponse.Response,
			Messages: errorResponse.Messages,
		}
	default:
		return errorResponse
	}
}

// parseRate referenced from https://github.com/google/go-github/blob/master/github/github.go#L495
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
			rate.Reset = Timestamp{time.Unix(v, 0)}
		}
	}
	return rate
}

// checkRateLimitBeforeDo referenced from https://github.com/google/go-github/blob/master/github/github.go#L627
func (c *Client) checkRateLimitBeforeDo(req *http.Request) *RateLimitError {
	c.rateMu.Lock()
	rate := c.rateLimit
	c.rateMu.Unlock()
	if !rate.Reset.Time.IsZero() && rate.Remaining == 0 && time.Now().Before(rate.Reset.Time) {
		// Create a fake response.
		resp := &http.Response{
			Status:     http.StatusText(http.StatusForbidden),
			StatusCode: http.StatusForbidden,
			Request:    req,
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(strings.NewReader("")),
		}
		return &RateLimitError{
			Rate:     rate,
			Response: resp,
			Messages: []string{fmt.Sprintf("API rate limit of %v still exceeded until %v, not making remote request.", rate.Limit, rate.Reset.Time)},
		}
	}

	return nil
}

// Rate referenced from https://github.com/google/go-github/blob/master/github/github.go#L861
type Rate struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     Timestamp `json:"reset"`
	err       error
}

// RateLimitError referenced from https://github.com/google/go-github/blob/master/github/github.go#L687
type RateLimitError struct {
	Rate     Rate           // Rate specifies last known rate limit for the client
	Response *http.Response // HTTP response that caused this error
	Messages []string       `json:"message"` // error message
}

type Meta struct {
	PreviousPage string `json:"previous_page"`
	NextPage     string `json:"next_page"`
	Total        int    `json:"total"`
}

// Error referenced from https://github.com/google/go-github/blob/master/github/github.go#L693
func (r *RateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Messages, r.Rate.Reset.Time.Sub(time.Now()))
}

// sanitizeURL referenced from https://github.com/google/go-github/blob/master/github/github.go#L734
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}

// ErrorResponse referenced from https://github.com/google/go-github/blob/master/github/github.go#L655
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Messages []string       `json:"messages"` // error message
	ErrorStr string         `json:"error"`    // more detail about an error
	err      error          // CheckResponse error
}

// ErrorResponse referenced from https://github.com/google/go-github/blob/master/github/github.go#L655
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Messages, r.ErrorStr)
}
