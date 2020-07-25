package docbase

import (
	"fmt"
	"net/http"
	"net/url"
)

//POST /teams/:domain/posts/:id/comments

type CommentService struct {
	client *Client
}

func NewCommentService(client *Client) *CommentService {
	return &CommentService{
		client: client,
	}
}

type Comment struct {
	ID int
	Body string
	Notice bool
//	author_id
//	published_at
}

type CommentRequest struct {
	Body  string `json:"body"`
}

type CommentResponse struct {
	Body  string `json:"body"`
}


func (s *CommentService) Create(postID string, cReq *CommentRequest) (*CommentResponse, *http.Response, error) {

	u, err := url.Parse(fmt.Sprintf("/posts/%s/comments", postID))

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u.String(), cReq)

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


func (s *CommentService) Delete(commentID string) (*http.Response, error) {
	u, err := url.Parse(fmt.Sprintf("/comments/%s", commentID))

	if err != nil {

	}

	req, err := s.client.NewRequest("DELETE", u.String(), nil)

	if err != nil {

	}

	cResp := &CommentResponse{}
	resp, err := s.client.Do(req, cResp)
	if err != nil {
		return resp, err
	}

	return resp, err
}
