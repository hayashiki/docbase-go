package docbase

import (
	"net/http"
	"net/url"
)

// TagService implements interface with API /tags endpoint.
// See https://help.docbase.io/posts/45703#%E3%82%BF%E3%82%B0
type TagService interface {
	List() (*TagListResponse, *Response, error)
}

// tagService handles communication with API
type tagService struct {
	client *Client
}

// Tag represents a docbase Tag
type Tag struct {
	Name string `json:"name"`
}

type TagListResponse []Tag

func (s *tagService) List() (*TagListResponse, *Response, error) {
	u, err := url.Parse("/tags")

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u.String(), nil)

	tagResp := &TagListResponse{}
	resp, err := s.client.Do(req, tagResp)
	if err != nil {
		return nil, resp, err
	}

	return tagResp, resp, err
}
