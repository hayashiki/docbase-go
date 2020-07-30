package docbase

import (
	"net/http"
	"net/url"
)

type TagService struct {
	client *Client
}

type Tag struct {
	Name string `json:"name"`
}

type TagListResponse []Tag

func NewTagService(client *Client) *TagService {
	return &TagService{client: client}
}

func (s *TagService) List() (*TagListResponse, *http.Response, error) {
	u, err := url.Parse("/tags")

	req, err := s.client.NewRequest(http.MethodGet, u.String(), nil)

	tagResp := &TagListResponse{}
	resp, err := s.client.Do(req, tagResp)
	if err != nil {
		return nil, resp, err
	}

	return tagResp, resp, err
}
