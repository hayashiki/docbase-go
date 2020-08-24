package docbase

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// UserService implements interface with API /posts endpoint.
// See https://help.docbase.io/posts/45703#%E3%83%81%E3%83%BC%E3%83%A0
type UserService interface {
	List(opts *UserListOptions) (*UserListResponse, *Response, error)
}

// UserCli handles communication with API
type UserCli struct {
	client *Client
}

// User represents a docbase User
type User struct {
	ID                    int           `json:"id"`
	Name                  string        `json:"name"`
	Username              string        `json:"username"`
	ProfileImageURL       string        `json:"profile_image_url"`
	Role                  string        `json:"role"`
	PostsCount            int           `json:"posts_count"`
	LastAccessTime        time.Time     `json:"last_access_time"`
	TwoStepAuthentication bool          `json:"two_step_authentication"`
	Groups                []SimpleGroup `json:"groups"`
}

type SimpleUser struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	ProfileImageURL string `json:"profile_image_url"`
}

type UserListResponse []User

// UserListOptions identifies as query params of User List request
type UserListOptions struct {
	Q                 string `url:"q,omitempty"`
	Page              int    `url:"page,omitempty"`
	PerPage           int    `url:"per_page,omitempty"`
	IncludeUserGroups bool   `url:"include_user_groups,omitempty"`
}

// List User
func (s *UserCli) List(opts *UserListOptions) (*UserListResponse, *Response, error) {
	u, err := url.Parse("/users")

	if err != nil {
		return nil, nil, err
	}

	q := u.Query()
	q.Set("per_page", strconv.Itoa(opts.PerPage))
	q.Set("page", strconv.Itoa(opts.Page))
	q.Set("q", opts.Q)
	q.Set("include_user_groups", opts.Q)
	u.RawQuery = q.Encode()

	req, err := s.client.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, nil, err
	}

	userResp := &UserListResponse{}
	resp, err := s.client.Do(req, userResp)
	if err != nil {
		return nil, resp, err
	}

	return userResp, resp, err
}
