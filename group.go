package docbase

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type GroupService struct {
	client *Client
}

func NewGroupService(client *Client) *GroupService {
	return &GroupService{client: client}
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

type Group struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	PostsCount     int       `json:"posts_count"`
	LastActivityAt time.Time `json:"last_activity_at"`
	CreatedAt      time.Time `json:"created_at"`
	Users          []struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"users"`
}

func (s *GroupService) List() {

}

func (s *GroupService) Get(id string) (*Group, *http.Response, error) {
	u, err := url.Parse(fmt.Sprintf("/groups/%s", id))

	//	TODO: return if not have scope permission

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, nil, err
	}

	var res Group
	resp, err := s.client.Client.Do(req)

	if err != nil {
		return nil, nil, err
	}

	return &res, resp, err
}
