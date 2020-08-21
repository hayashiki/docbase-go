package docbase

import (
	"fmt"
	"github.com/hayashiki/docbase-go/testutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestAttachmentService_Download(t *testing.T) {
	setup()
	defer teardown()

	ti, err := time.Parse(time.RFC3339, "2020-03-27T09:25:09+09:00")

	att1 := Attachment{
		ID:        "fd26b8c9-0c55-48e7-a943-87292acd5682.png",
		Name:      "image1.png",
		Size:      132323,
		URL:       "https://image.docbase.io/uploads/fd26b8c9-0c55-48e7-a943-87292acd5682.png",
		Markdown:  "![image.png](https://image.docbase.io/uploads/fd26b8c9-0c55-48e7-a943-87292acd5682.png)",
		CreatedAt: ti,
	}

	mux.HandleFunc(fmt.Sprintf("/attachments/%s", att1.ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{}`)
	})

	attSvc := AttachmentCli{client}

	_, resp, err := attSvc.Download("fd26b8c9-0c55-48e7-a943-87292acd5682.png")

	if err != nil {
		t.Errorf("Shouldn't have returned an error: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Download attachment response code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}

	//TODO
	//if !reflect.DeepEqual(attachment, att1) {
	//	t.Errorf("Group returned %+v, want %+v", attachment, att1)
	//}
}

func TestAttachmentService_Upload(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/attachments", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		testMethod(t, r, "POST")

		fmt.Fprint(w, testutil.LoadFixture(t, "attachment-list-response.json"))
	})

	attSvc := AttachmentCli{client}

	atts, resp, err := attSvc.Upload([]string{"./testdata/image1.jpg", "./testdata/image2.jpg"})

	if err != nil {
		t.Errorf("Shouldn't have returned an error: %+v", err)
	}

	ti, err := time.Parse(time.RFC3339, "2020-03-27T09:25:09+09:00")

	att1 := Attachment{
		ID:        "fd26b8c9-0c55-48e7-a943-87292acd5682.png",
		Name:      "image1.png",
		Size:      132323,
		URL:       "https://image.docbase.io/uploads/fd26b8c9-0c55-48e7-a943-87292acd5682.png",
		Markdown:  "![image.png](https://image.docbase.io/uploads/fd26b8c9-0c55-48e7-a943-87292acd5682.png)",
		CreatedAt: ti,
	}

	att2 := Attachment{
		ID:        "gd26b8c9-0c55-48e7-a943-87292acd5683.png",
		Name:      "image2.png",
		Size:      132324,
		URL:       "https://image.docbase.io/uploads/gd26b8c9-0c55-48e7-a943-87292acd5683.png",
		Markdown:  "![image.png](https://image.docbase.io/uploads/gd26b8c9-0c55-48e7-a943-87292acd5683.png)",
		CreatedAt: ti,
	}

	want := &AttachmentResponse{att1, att2}

	if err != nil {
		t.Errorf("Fail to get group request err: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Upload attachment response code = %v, expected %v", resp.StatusCode, http.StatusCreated)
	}

	if !reflect.DeepEqual(atts, want) {
		t.Errorf("Attachment returned %+v, want %+v", atts, want)
	}
}
