package docbase

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// CommentService implements interface with API /groups endpoint.
// https://help.docbase.io/posts/45703#%E3%82%B3%E3%83%A1%E3%83%B3%E3%83%88
type CommentService interface {
	Create(postID int, commentRequest *CommentCreateRequest) (*Comment, *Response, error)
	Delete(commentID int) (*Response, error)
}

// CommentCli handles communication with API
type CommentCli struct {
	client *Client
}

// Comment represents a docbase Comment
type Comment struct {
	ID         int       `json:"id"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
	SimpleUser `json:"user"`
}

// CommentCreateRequest identifies Comment for the Create request
type CommentCreateRequest struct {
	Body        string    `json:"body"`
	Notice      bool      `json:"notice,omitempty"`
	AuthorID    string    `json:"author_id,omitempty"`
	PublishedAt time.Time `json:"published_at,omitempty"`
}

// Create Comment
func (s *CommentCli) Create(postID int, commentRequest *CommentCreateRequest) (*Comment, *Response, error) {

	u, err := url.Parse(fmt.Sprintf("/posts/%d/comments", postID))

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u.String(), commentRequest)

	if err != nil {
		return nil, nil, err
	}

	cResp := &Comment{}
	resp, err := s.client.Do(req, cResp)
	if err != nil {
		return nil, resp, err
	}

	return cResp, resp, err
}

// Delete Comment
func (s *CommentCli) Delete(commentID int) (*Response, error) {
	u, err := url.Parse(fmt.Sprintf("/comments/%d", commentID))

	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodDelete, u.String(), nil)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
