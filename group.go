package docbase

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// GroupService implements interface with API /groups endpoint.
// See https://help.docbase.io/posts/45703#%E3%82%B0%E3%83%AB%E3%83%BC%E3%83%97
type GroupService interface {
	List(opts *GroupListOptions) (*GroupListResponse, *Response, error)
	Get(id int) (*Group, *Response, error)
	Create(createRequest *GroupCreateRequest) (*Group, *Response, error)
	// TODO: create group users file
	AddUser(id int, gReq *GroupAddRequest) (*Response, error)
	RemoveUser(id int, gReq *GroupAddRequest) (*Response, error)
}

// GroupCli handles communication with API
type GroupCli struct {
	client *Client
}

// Group represents a minimum group Object
type SimpleGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Group represents a group Object
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

// GroupListOptions identifies as query params of List request
type GroupListOptions struct {
	Name    string `url:"name,omitempty"`
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
}

type GroupCreateRequest struct {
	Name  string `json:"name"`
	Description   string `json:"description"`
	//Draft  bool   `json:"draft"`  // optional, default: false
	//Notice bool   `json:"notice"` // optional, default: true
	//Tags   []string
	//Scope  string `json:"scope"` // optional, default: everyone
	//Groups []string
}

// GroupListResponse represents a List simple group
type GroupListResponse []SimpleGroup

// List Group
func (s *GroupCli) List(opts *GroupListOptions) (*GroupListResponse, *Response, error) {
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
		return nil, resp, err
	}

	return res, resp, err
}

// Get Group
func (s *GroupCli) Get(id int) (*Group, *Response, error) {
	u, err := url.Parse(fmt.Sprintf("/groups/%d", id))

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

// Create Group
func (s *GroupCli) Create(createRequest *GroupCreateRequest) (*Group, *Response, error) {
	u, err := url.Parse("/groups")

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u.String(), createRequest)

	if err != nil {
		return nil, nil, err
	}

	cResp := &Group{}
	resp, err := s.client.Do(req, cResp)
	if err != nil {
		return nil, resp, err
	}

	return cResp, resp, err
}

func (s *GroupCli) AddUser(id int, gReq *GroupAddRequest) (*Response, error) {
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

func (s *GroupCli) RemoveUser(id int, gReq *GroupAddRequest) (*Response, error) {
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

func NewGroupService(client *Client) *GroupCli {
	return &GroupCli{client: client}
}
