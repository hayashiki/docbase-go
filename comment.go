package docbase

import (
	"fmt"
	"net/http"
	"net/url"
)

// CommentService implements interface with API /groups endpoint.
// https://help.docbase.io/posts/45703#%E3%82%B3%E3%83%A1%E3%83%B3%E3%83%88
type CommentService interface {
	Create(postID int, commentRequest *CommentRequest) (*CommentResponse, *Response, error)
	Delete(commentID int) (*Response, error)
}

type CommentCli struct {
	client *Client
}

func NewCommentService(client *Client) *CommentCli {
	return &CommentCli{
		client: client,
	}
}

type Comment struct {
	ID     int
	Body   string
	Notice bool
	//	author_id
	//	published_at
}

type CommentRequest struct {
	Body string `json:"body"`
}

type CommentResponse struct {
	Body string `json:"body"`
}

func (s *CommentCli) Create(postID int, cReq *CommentRequest) (*CommentResponse, *Response, error) {

	u, err := url.Parse(fmt.Sprintf("/posts/%d/comments", postID))

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u.String(), cReq)

	if err != nil {
		return nil, nil, err
	}

	cResp := &CommentResponse{}
	resp, err := s.client.Do(req, cResp)
	if err != nil {
		return nil, resp, err
	}

	return cResp, resp, err
}

func (s *CommentCli) Delete(commentID int) (*Response, error) {
	u, err := url.Parse(fmt.Sprintf("/comments/%d", commentID))

	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodDelete, u.String(), nil)

	if err != nil {
		return nil, err
	}

	cResp := &CommentResponse{}
	resp, err := s.client.Do(req, cResp)
	if err != nil {
		return resp, err
	}

	return resp, err
}
