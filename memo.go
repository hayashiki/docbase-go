package docbase

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type MemoService struct {
	client *Client
}

func NewMemoService(client *Client) *MemoService {
	return &MemoService{client: client}
}

type MemoRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Draft  bool   `json:"draft"`  // optional, default: false
	Notice bool   `json:"notice"` // optional, default: true
	Tags   []string
	Scope  string `json:"scope"` // optional, default: everyone
	Groups []string
}

type MemoCreateResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Draft     bool      `json:"draft"`
	Archived  bool      `json:"archived"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []struct {
		Name string `json:"name"`
	} `json:"tags"`
	Scope      string `json:"scope"`
	SharingURL string `json:"sharing_url"`
	User       struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"user"`
	StarsCount    int           `json:"stars_count"`
	GoodJobsCount int           `json:"good_jobs_count"`
	Comments      []interface{} `json:"comments"`
	Groups        []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"groups"`
	Attachments []struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Size      int       `json:"size"`
		URL       string    `json:"url"`
		Markdown  string    `json:"markdown"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"attachments"`
}

func (s *MemoService) Create(memoReq *MemoRequest) (*MemoCreateResponse, *http.Response, error) {
	u, err := url.Parse("/posts")

	//	TODO: return if not have scope permission

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u.String(), memoReq)

	if err != nil {
		return nil, nil, err
	}

	var res MemoCreateResponse
	resp, err := s.client.Do(req, res)

	if err != nil {
		return nil, nil, err
	}

	return &res, resp, err
}

func (s *MemoService) Get(memoID string) (*MemoCreateResponse, *http.Response, error) {
	u, err := url.Parse(fmt.Sprintf("/posts/%s", memoID))

	//	TODO: return if not have scope permission

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, nil, err
	}

	var res MemoCreateResponse
	resp, err := s.client.Client.Do(req)

	if err != nil {
		return nil, nil, err
	}

	return &res, resp, err
}
