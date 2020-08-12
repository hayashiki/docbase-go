package docbase

import (
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

//curl \
//  -H 'X-DocBaseToken: ACCESS TOKEN' \
//  https://api.docbase.io/teams/kray/users?include_user_groups=true
//q	ユーザ名もしくはユーザIDの一部	String
//page	ページ番号	Integer		1
//per_page	1ページのユーザ数	Integer		100	100
//include_user_groups ※	ユーザの所属グループを含めるかどうか	Boolean		false

type UserListOptions struct {
	Q       string
	Page    int
	PerPage int
}

func (s *UserCli) List(opts *UserListOptions) (*UserListResponse, *Response, error) {
	u, err := url.Parse("/users")

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

	userResp := &UserListResponse{}
	resp, err := s.client.Do(req, userResp)
	if err != nil {
		return nil, resp, err
	}

	return userResp, resp, err
}

func NewUserService(client *Client) *UserCli {
	return &UserCli{client: client}
}
