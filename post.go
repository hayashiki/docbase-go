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
	Create(postRequest *PostCreateRequest) (*Post, *Response, error)
	Update(postID int, postUpdateRequest *PostUpdateRequest) (*Post, *Response, error)
	Delete(postID string) (*Response, error)
	Archive(postID int) (*Response, error)
	Unarchive(postID int) (*Response, error)
}

// PostCli handles communication with API
type PostCli struct {
	client *Client
}

// PostCreateRequest identifies Post for the Create request
type PostCreateRequest struct {
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Draft       bool      `json:"draft"`  // optional, default: false
	Notice      bool      `json:"notice"` // optional, default: true
	Tags        []string  `json:"tags"`
	Scope       string    `json:"scope"` // optional, default: everyone
	Groups      []string  `json:"groups"`
	AuthorID    string    `json:"author_id"`
	PublishedAt time.Time `json:"published_at"`
}

// PostUpdateRequest identifies Post for the Update request
type PostUpdateRequest struct {
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Draft       bool      `json:"draft"`  // optional, default: false
	Notice      bool      `json:"notice"` // optional, default: true
	Tags        []string  `json:"tags"`
	Scope       string    `json:"scope"` // optional, default: everyone
	Groups      []string  `json:"scope"`
	PublishedAt time.Time `json:"published_at"`
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
	ID            int           `json:"id"`
	Title         string        `json:"title"`
	Body          string        `json:"body"`
	Draft         bool          `json:"draft"`
	Archived      bool          `json:"archived"`
	URL           string        `json:"url"`
	CreatedAt     time.Time     `json:"created_at"`
	Tags          []Tag         `json:"tags"`
	Scope         string        `json:"scope"`
	SharingURL    string        `json:"sharing_url"`
	User          SimpleUser    `json:"user"`
	StarsCount    int           `json:"stars_count"`
	GoodJobsCount int           `json:"good_jobs_count"`
	Comments      []Comment     `json:"comments"`
	Groups        []SimpleGroup `json:"groups"`
	Attachments   []Attachment  `json:"attachments"`
}

// PostListOptions identifies as query params of Post List request
type PostListOptions struct {
	Q       string `url:"q,omitempty"`
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
}

// List Post
func (s *PostCli) List(opts *PostListOptions) (*PostListResponse, *Response, error) {

	u, err := url.Parse("/posts")

	if err != nil {
		return nil, nil, err
	}

	q := u.Query()
	q.Set("per_page", strconv.Itoa(opts.PerPage))
	q.Set("page", strconv.Itoa(opts.Page))
	q.Set("q", opts.Q)

	req, err := s.client.NewRequest(http.MethodGet, u.String(), nil)

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

// Get Post
func (s *PostCli) Get(postID int) (*Post, *Response, error) {

	u, err := url.Parse(fmt.Sprintf("/posts/%d", postID))

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u.String(), nil)

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

// Create Post
func (s *PostCli) Create(memoReq *PostCreateRequest) (*Post, *Response, error) {
	u, err := url.Parse("/posts")

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u.String(), memoReq)

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

// Update Post
func (s *PostCli) Update(postID int, postUpdateRequest *PostUpdateRequest) (*Post, *Response, error) {
	u, err := url.Parse(fmt.Sprintf("/posts/%d", postID))
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPatch, u.String(), postUpdateRequest)

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

// Delete Post
func (s *PostCli) Delete(postID string) (*Response, error) {
	u, err := url.Parse(fmt.Sprintf("/posts/%s", postID))
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

// Archive Post
func (s *PostCli) Archive(postID int) (*Response, error) {
	u, err := url.Parse(fmt.Sprintf("/posts/%d/archive", postID))
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

// Unarchive Post
func (s *PostCli) Unarchive(postID int) (*Response, error) {
	u, err := url.Parse(fmt.Sprintf("/posts/%d/unarchive", postID))
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
