package docbase

import (
	"fmt"
	"github.com/hayashiki/docbase-go/testutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestMemoService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {

		log.Printf("posts is")

		fmt.Fprintf(w, testutil.LoadFixture(t,"memo-detail-response.json"))
	})

	memoSrv := NewMemoService(client)

	mReq := &MemoRequest{}

	actual, _, err := memoSrv.Create(mReq)

	if err != nil {
		t.Errorf("Shouldn't have returned an error: %+v", err)
	}

	ti, err := time.Parse(time.RFC3339, "2020-03-27T09:25:09+09:00")

	if err != nil {
		t.Errorf("Fail to parse err: %v", err)
	}

	want := &Memo{
		ID:         1,
		Title:      "メモのタイトル",
		Body:       "メモの本文",
		Draft:      false,
		Archived:   false,
		URL:        "https://kray.docbase.io/posts/1",
		CreatedAt:  ti,
		Tags:       []MemoTag{
			MemoTag{Name: "rails"},
			MemoTag{Name: "ruby"},
		},
		Scope:      "group",
		SharingURL: "https://docbase.io/posts/1/sharing/abcdefgh-0e81-4567-9876-1234567890ab",
		User: MemoUser{
			ID:              1,
			Name:            "danny",
			ProfileImageURL: "https://image.docbase.io/uploads/aaa.gif",
		},
		StarsCount:    1,
		GoodJobsCount: 2,
		Comments:      []MemoComment{},
		Groups:        []MemoGroup{
			MemoGroup{
				ID: 1,
				Name: "DocBase",
			},
		},
		Attachments:   []MemoAttachment{MemoAttachment{
			ID:        "461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv",
			Name:      "uploadfile.csv",
			Size:      18786,
			URL:       "https://kray.docbase.io/file_attachments/461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv",
			Markdown:  "[![csv](/images/file_icons/csv.svg) uploadfile.jpg](https://kray.docbase.io/uploads/461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv)",
			CreatedAt: ti,
		},
		},
	}

	if !reflect.DeepEqual(actual, want) {
		t.Errorf("Get return %+v, want %+v", actual, want)
	}
}

func TestMemoService_Get(t *testing.T) {

	setup()
	defer teardown()

	ti, err := time.Parse(time.RFC3339, "2020-03-27T09:25:09+09:00")

	if err != nil {
		t.Errorf("Fail to parse err: %v", err)
	}

	memo := &Memo{
		ID:         1,
		Title:      "メモのタイトル",
		Body:       "メモの本文",
		Draft:      false,
		Archived:   false,
		URL:        "https://kray.docbase.io/posts/1",
		CreatedAt:  ti,
		Tags:       []MemoTag{
			MemoTag{Name: "rails"},
			MemoTag{Name: "ruby"},
		},
		Scope:      "group",
		SharingURL: "https://docbase.io/posts/1/sharing/abcdefgh-0e81-4567-9876-1234567890ab",
		User: MemoUser{
			ID:              1,
			Name:            "danny",
			ProfileImageURL: "https://image.docbase.io/uploads/aaa.gif",
		},
		StarsCount:    1,
		GoodJobsCount: 2,
		Comments:      []MemoComment{},
		Groups:        []MemoGroup{
			MemoGroup{
				ID: 1,
				Name: "DocBase",
			},
		},
		Attachments:   []MemoAttachment{MemoAttachment{
			ID:        "461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv",
			Name:      "uploadfile.csv",
			Size:      18786,
			URL:       "https://kray.docbase.io/file_attachments/461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv",
			Markdown:  "[![csv](/images/file_icons/csv.svg) uploadfile.jpg](https://kray.docbase.io/uploads/461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv)",
			CreatedAt: ti,
		},
		},
	}

	memoSvc := NewMemoService(client)

	mux.HandleFunc(fmt.Sprintf("/posts/%d", memo.ID), func(w http.ResponseWriter, r *http.Request) {
		//requestSent = true

		u, _ := url.Parse(fmt.Sprintf("/posts/%d", memo.ID))

		want := u.String()
		if got := r.URL.String(); got != want {
			t.Errorf("URL: got %v, want %v", got, want)
		}

		fmt.Fprint(w, testutil.LoadFixture(t, "memo-detail-response.json"))
	})

	getRes, _, err := memoSvc.Get(strconv.Itoa(memo.ID))

	if err != nil {
		t.Errorf("Get returned an error: %v", err)
	}

	if !reflect.DeepEqual(getRes, memo) {
		t.Errorf("Get returned %+v, want %+v", getRes, memo)
	}
}

