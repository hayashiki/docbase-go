package docbase

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (f *File) Encode(filePath string) error {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return err
	}

	f.Name = file.Name()
	fi, _ := file.Stat()
	size := fi.Size()
	data := make([]byte, size)
	_, err = file.Read(data)

	if err != nil {
		return err
	}

	f.Content = base64.StdEncoding.EncodeToString(data)

	return nil
}

type AttachmentService struct {
	client *Client
}

type Attachment struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Size      int       `json:"size"`
	URL       string    `json:"url"`
	Markdown  string    `json:"markdown"`
	CreatedAt time.Time `json:"created_at"`
}

type AttachmentResponse []Attachment

type FileContent []byte

type Request struct {
	ID string
}

func (s *AttachmentService) Download(attachmentID string) (*FileContent, *Response, error) {
	u, err := url.Parse(fmt.Sprintf("/attachments/%s", attachmentID))

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, nil, err
	}

	fileResp, resp, err := s.client.DoUpload(req)

	if err != nil {
		return nil, resp, err
	}

	return &fileResp, resp, nil
}

func (s *AttachmentService) Upload(filesPath []string) (*AttachmentResponse, *Response, error) {

	var files []File

	for _, fp := range filesPath {
		var file File
		err := file.Encode(fp)

		if err != nil {
			return nil, nil, fmt.Errorf("failed read file err: %w", err)
		}

		files = append(files, file)
	}

	u, err := url.Parse("/attachments")

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u.String(), files)

	if err != nil {
		return nil, nil, err
	}

	atRes := &AttachmentResponse{}
	resp, err := s.client.Do(req, atRes)

	if err != nil {
		return nil, resp, err
	}

	return atRes, resp, nil
}

func NewAttachmentService(client *Client) *AttachmentService {
	return &AttachmentService{client: client}
}
