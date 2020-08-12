package docbase

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// PostService implements interface with API /posts endpoint.
// See https://help.docbase.io/posts/45703#%E3%82%BF%E3%82%B0
type PostService interface {
	List(opts *PostListOptions) (*PostListResponse, *Response, error)
	Get(postID int) (*Post, *Response, error)
	Create(postRequest *PostRequest) (*Post, *Response, error)
	Delete(postID string) (*Response, error)
	Update(postID int, postRequest *PostRequest) (*Post, *Response, error)
	Archive(postID int) (*Response, error)
	Unarchive(postID int) (*Response, error)
}

// PostCli handles communication with API
type PostCli struct {
	client *Client
}

type PostRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Draft  bool   `json:"draft"`  // optional, default: false
	Notice bool   `json:"notice"` // optional, default: true
	Tags   []string
	Scope  string `json:"scope"` // optional, default: everyone
	Groups []string
}

type PostListResponse struct {
	Posts []Post `json:"posts"`
	Meta  struct {
		PreviousPage string `json:"previous_page"`
		NextPage     string `json:"next_page"`
		Total        int    `json:"total"`
	} `json:"meta"`
}

// Post represents a docbase Post
type Post struct {
	ID            int              `json:"id"`
	Title         string           `json:"title"`
	Body          string           `json:"body"`
	Draft         bool             `json:"draft"`
	Archived      bool             `json:"archived"`
	URL           string           `json:"url"`
	CreatedAt     time.Time        `json:"created_at"`
	Tags          []Tag            `json:"tags"`
	Scope         string           `json:"scope"`
	SharingURL    string           `json:"sharing_url"`
	User          SimpleUser       `json:"user"`
	StarsCount    int              `json:"stars_count"`
	GoodJobsCount int              `json:"good_jobs_count"`
	Comments      []PostComment    `json:"comments"`
	Groups        []SimpleGroup    `json:"groups"`
	Attachments   []PostAttachment `json:"attachments"`
}

type PostListOptions struct {
	Q       string
	Page    int
	PerPage int
}

type PostComment struct {
	ID        int
	Body      string
	CreatedAt time.Time `json:"created_at"`
	SimpleUser
}

type PostAttachment struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Size      int       `json:"size"`
	URL       string    `json:"url"`
	Markdown  string    `json:"markdown"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *PostCli) List(opts *PostListOptions) (*PostListResponse, *Response, error) {

	u, err := url.Parse("/posts")

	if err != nil {
		return nil, nil, err
	}

	q := u.Query()
	q.Set("per_page", strconv.Itoa(opts.PerPage))
	q.Set("page", strconv.Itoa(opts.Page))
	q.Set("q", opts.Q)

	req, err := s.client.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, nil, err
	}

	mResp := &PostListResponse{}
	resp, err := s.client.Do(req, mResp)

	if err != nil {
		return nil, nil, err
	}

	return mResp, resp, err
}

func (s *PostCli) Get(postID int) (*Post, *Response, error) {

	u, err := url.Parse(fmt.Sprintf("/posts/%d", postID))

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, nil, err
	}

	mResp := &Post{}
	resp, err := s.client.Do(req, mResp)

	if err != nil {
		return nil, nil, err
	}

	return mResp, resp, err
}

func (s *PostCli) Create(memoReq *PostRequest) (*Post, *Response, error) {
	u, err := url.Parse("/posts")

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u.String(), memoReq)

	if err != nil {
		return nil, nil, err
	}

	mResp := &Post{}
	resp, err := s.client.Do(req, mResp)

	if err != nil {
		return nil, nil, err
	}

	return mResp, resp, err
}

func (s *PostCli) Update(memoID int, memoReq *PostRequest) (*Post, *Response, error) {
	u, err := url.Parse(fmt.Sprintf("/posts/%d", memoID))
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPatch, u.String(), memoReq)

	if err != nil {
		return nil, nil, err
	}

	mResp := &Post{}
	resp, err := s.client.Do(req, mResp)
	if err != nil {
		return nil, nil, err
	}

	if err != nil {
		return nil, nil, err
	}

	return mResp, resp, err
}

func (s *PostCli) Delete(memoID string) (*Response, error) {
	u, err := url.Parse(fmt.Sprintf("/posts/%s", memoID))
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodDelete, u.String(), nil)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *PostCli) Archive(memoID int) (*Response, error) {
	u, err := url.Parse(fmt.Sprintf("/posts/%d/archive", memoID))
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPut, u.String(), nil)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return resp, err
}

func (s *PostCli) Unarchive(memoID int) (*Response, error) {
	u, err := url.Parse(fmt.Sprintf("/posts/%d/unarchive", memoID))
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPut, u.String(), nil)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return resp, err
}

func NewPostService(client *Client) *PostCli {
	return &PostCli{client: client}
}
