package docbase

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type AttachmentService struct {
	client *Client
}

type Attachment struct {
	ID        string       `json:"id"`
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

func (s *AttachmentService) Download(attachmentID string) (*FileContent, *http.Response, error) {
	u, err := url.Parse(fmt.Sprintf("/attachments/%s", attachmentID))

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u.String(), nil)

	if err != nil {
		log.Print("hoge2")
		return nil, nil, err
	}

	fileResp, resp, err := s.client.DoBinary(req)

	if err != nil {
		log.Print("hoge3")
		return nil, resp, err
	}

	return &fileResp, resp, nil
}

type File struct {
	Name string `json:"name"`
	Content string `json:"content"`
}

func (f *File) Encode(filePath string) {
	file, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	f.Name = file.Name()
	fi, _ := file.Stat()
	size := fi.Size()
	data := make([]byte, size)
	file.Read(data)

	f.Content = base64.StdEncoding.EncodeToString(data)
}

func (s *AttachmentService) Upload(filePath string) (*AttachmentResponse, *http.Response, error) {

	var files []File
	var file File
	file.Encode(filePath)

	files = append(files, file)

	u, err := url.Parse(fmt.Sprintf("/attachments"))

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

