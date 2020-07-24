package docbase

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type UserService struct {
	client *Client
}

type User struct {
	ID                    int       `json:"id"`
	Name                  string    `json:"name"`
	Username              string    `json:"username"`
	ProfileImageURL       string    `json:"profile_image_url"`
	Role                  string    `json:"role"`
	PostsCount            int       `json:"posts_count"`
	LastAccessTime        time.Time `json:"last_access_time"`
	TwoStepAuthentication bool      `json:"two_step_authentication"`
	Groups                []UGroup  `json:"groups"`
}

type UGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//curl \
//  -H 'X-DocBaseToken: ACCESS TOKEN' \
//  https://api.docbase.io/teams/kray/users?include_user_groups=true
//q	ユーザ名もしくはユーザIDの一部	String
//page	ページ番号	Integer		1
//per_page	1ページのユーザ数	Integer		100	100
//include_user_groups ※	ユーザの所属グループを含めるかどうか	Boolean		false

type UserListOptions struct {
	Q string
	Page int
	PerPage int
}

func NewUserService(client *Client) *UserService {
	return &UserService{client: client}
}

func (s *UserService) List(opts *UserListOptions) (*[]User, *http.Response, error) {
	u, err := url.Parse("/users")

	if err != nil {

	}

	q := u.Query()

	q.Set("per_page", strconv.Itoa(opts.PerPage))
	q.Set("page", strconv.Itoa(opts.Page))
	q.Set("q", opts.Q)

	req, err := s.client.NewRequest("GET", u.String(), nil)

	if err != nil {

	}

	userResp := &[]User{}
	resp, err := s.client.Do(req, userResp)
	if err != nil {
		return nil, resp, err
	}

	return userResp, resp, err
}
