package docbase

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type GroupService struct {
	client *Client
}

type SimpleGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Group struct {
	ID             int          `json:"id"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	PostsCount     int          `json:"posts_count"`
	LastActivityAt time.Time    `json:"last_activity_at"`
	CreatedAt      time.Time    `json:"created_at"`
	Users          []SimpleUser `json:"users"`
}

type GroupAddRequest struct {
	UserIDs []int `json:"user_ids"`
}

type GroupListOptions struct {
	Name    string
	Page    int
	PerPage int
}

type GroupRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Draft  bool   `json:"draft"`  // optional, default: false
	Notice bool   `json:"notice"` // optional, default: true
	Tags   []string
	Scope  string `json:"scope"` // optional, default: everyone
	Groups []string
}

type GroupListResponse []SimpleGroup

func (s *GroupService) List(opts *GroupListOptions) (*GroupListResponse, *http.Response, error) {
	u, err := url.Parse("/groups")

	if err != nil {
		return nil, nil, err
	}

	q := u.Query()
	q.Set("per_page", strconv.Itoa(opts.PerPage))
	q.Set("page", strconv.Itoa(opts.Page))
	q.Set("q", opts.Name)

	req, err := s.client.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, nil, err
	}

	res := &GroupListResponse{}
	resp, err := s.client.Do(req, res)

	if err != nil {
		return nil, nil, err
	}

	return res, resp, err

}

func (s *GroupService) Get(id int) (*Group, *http.Response, error) {
	u, err := url.Parse(fmt.Sprintf("/groups/%d", id))

	//	TODO: return if not have scope permission

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, nil, err
	}

	res := &Group{}
	resp, err := s.client.Do(req, res)

	if err != nil {
		return nil, nil, err
	}

	return res, resp, err
}

func (s *GroupService) AddUser(id int, gReq *GroupAddRequest) (*http.Response, error) {
	u, err := url.Parse(fmt.Sprintf("/groups/%d/users", id))

	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u.String(), gReq)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func (s *GroupService) RemoveUser(id int, gReq *GroupAddRequest) (*http.Response, error) {
	u, err := url.Parse(fmt.Sprintf("/groups/%d/users", id))

	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodDelete, u.String(), gReq)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func NewGroupService(client *Client) *GroupService {
	return &GroupService{client: client}
}
