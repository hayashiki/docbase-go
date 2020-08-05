package docbase

import (
	"fmt"
	"net/http"
	"net/url"
)

type CommentService struct {
	client *Client
}

func NewCommentService(client *Client) *CommentService {
	return &CommentService{
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

func (s *CommentService) Create(postID int, cReq *CommentRequest) (*CommentResponse, *Response, error) {

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

func (s *CommentService) Delete(commentID int) (*Response, error) {
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
