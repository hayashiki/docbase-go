package docbase

import (
	"fmt"
	"net/http"
	"net/url"
)

// GroupUserService implements interface with API /groups/:id/users endpoint.
// See https://help.docbase.io/posts/45703#%E3%82%B0%E3%83%AB%E3%83%BC%E3%83%97
type GroupUserService interface {
	Create(id int, groupUserCreateRequest *GroupUserCreateRequest) (*Response, error)
	Delete(id int, groupUserCreateRequest *GroupUserCreateRequest) (*Response, error)
}

// groupUserService handles communication with API
type groupUserService struct {
	client *Client
}

type GroupUserCreateRequest struct {
	UserIDs []int `json:"user_ids"`
}

func (c *groupUserService) Create(id int, groupUserCreateRequest *GroupUserCreateRequest) (*Response, error) {
	u, err := url.Parse(fmt.Sprintf("/groups/%d/users", id))

	if err != nil {
		return nil, err
	}

	req, err := c.client.NewRequest(http.MethodPost, u.String(), groupUserCreateRequest)

	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req, nil)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func (c *groupUserService) Delete(id int, groupUserCreateRequest *GroupUserCreateRequest) (*Response, error) {
	u, err := url.Parse(fmt.Sprintf("/groups/%d/users", id))

	if err != nil {
		return nil, err
	}

	req, err := c.client.NewRequest(http.MethodDelete, u.String(), groupUserCreateRequest)

	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req, nil)

	if err != nil {
		return nil, err
	}

	return resp, err
}
