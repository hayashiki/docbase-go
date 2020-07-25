package docbase

import (
	"fmt"
	"log"
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

type Memo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Draft     bool      `json:"draft"`
	Archived  bool      `json:"archived"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []MemoTag `json:"tags"`
	Scope      string `json:"scope"`
	SharingURL string `json:"sharing_url"`
	User       MemoUser `json:"user"`
	StarsCount    int           `json:"stars_count"`
	GoodJobsCount int           `json:"good_jobs_count"`
	Comments      []MemoComment `json:"comments"`
	Groups        []MemoGroup `json:"groups"`
	Attachments []MemoAttachment `json:"attachments"`
}

type MemoComment struct {
	ID int
	Body string
	CreatedAt time.Time `json:"created_at"`
	MemoUser
}

type MemoTag struct {
	Name string `json:"name"`
}

type MemoUser struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	ProfileImageURL string `json:"profile_image_url"`
}

type MemoGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MemoAttachment struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Size      int       `json:"size"`
	URL       string    `json:"url"`
	Markdown  string    `json:"markdown"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *MemoService) Create(memoReq *MemoRequest) (*Memo, *http.Response, error) {
	u, err := url.Parse("/posts")

	//	TODO: return if not have scope permission

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u.String(), memoReq)

	if err != nil {
		return nil, nil, err
	}

	mResp := &Memo{}
	resp, err := s.client.Do(req, mResp)

	if err != nil {
		return nil, nil, err
	}

	log.Printf("res is %v", mResp )
	log.Printf("res is %v", err )

	return mResp, resp, err
}

func (s *MemoService) Get(memoID string) (*Memo, *http.Response, error) {

	u, err := url.Parse(fmt.Sprintf("/posts/%s", memoID))

	//	TODO: return if not have scope permission

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, nil, err
	}

	mResp := &Memo{}
	resp, err := s.client.Do(req, mResp)

	if err != nil {
		return nil, nil, err
	}

	return mResp, resp, err
}
