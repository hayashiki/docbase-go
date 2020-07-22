package docbase

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AttachmentService struct {
	client *Client
}

func NewAttachmentService(client *Client) *AttachmentService {
	return &AttachmentService{client: client}
}


type Attachment struct {

}

type Request struct {
	ID string
}

func (s AttachmentService) Download(attachmentID string) ([]byte, error){
	u, err := url.Parse(fmt.Sprintf("%s/attachments/%s", s.client.BaseURL, attachmentID))

	if err != nil {
		fmt.Errorf("image file is %v", err)
		return nil, nil
	}

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		fmt.Errorf("image file is %v", err)
		return nil, nil
	}

	req.Header.Add("X-DocBaseToken", s.client.AccessToken)

	resp, err := s.client.Client.Do(req)

	if err != nil {
		fmt.Errorf("image file is %v", err)
		return nil, nil
	}

	fileBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Errorf("image file is %v", err)
		return nil, nil
	}

	return fileBytes, nil
}
