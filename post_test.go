package docbase

import (
	"fmt"
	"github.com/hayashiki/docbase-go/testutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestPostService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		testMethod(t, r, "POST")

		fmt.Fprint(w, testutil.LoadFixture(t, "post-detail-response.json"))
	})

	mReq := &PostCreateRequest{}

	post, resp, err := client.Posts.Create(mReq)

	if err != nil {
		t.Errorf("Shouldn't have returned an error: %+v", err)
	}

	ti, err := time.Parse(time.RFC3339, "2020-03-27T09:25:09+09:00")

	if err != nil {
		t.Errorf("Fail to parse err: %v", err)
	}

	want := &Post{
		ID:        1,
		Title:     "メモのタイトル",
		Body:      "メモの本文",
		Draft:     false,
		Archived:  false,
		URL:       "https://kray.docbase.io/posts/1",
		CreatedAt: ti,
		Tags: []Tag{
			Tag{Name: "rails"},
			Tag{Name: "ruby"},
		},
		Scope:      "group",
		SharingURL: "https://docbase.io/posts/1/sharing/abcdefgh-0e81-4567-9876-1234567890ab",
		User: SimpleUser{
			ID:              1,
			Name:            "danny",
			ProfileImageURL: "https://image.docbase.io/uploads/aaa.gif",
		},
		StarsCount:    1,
		GoodJobsCount: 2,
		Comments:      []Comment{},
		Groups: []SimpleGroup{
			SimpleGroup{
				ID:   1,
				Name: "DocBase",
			},
		},
		Attachments: []Attachment{Attachment{
			ID:        "461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv",
			Name:      "uploadfile.csv",
			Size:      18786,
			URL:       "https://kray.docbase.io/file_attachments/461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv",
			Markdown:  "[![csv](/images/file_icons/csv.svg) uploadfile.jpg](https://kray.docbase.io/uploads/461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv)",
			CreatedAt: ti,
		},
		},
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Post Create request code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}

	if !reflect.DeepEqual(post, want) {
		t.Errorf("Get return %+v, want %+v", post, want)
	}
}

func TestPostService_Get(t *testing.T) {
	setup()
	defer teardown()

	ti, err := time.Parse(time.RFC3339, "2020-03-27T09:25:09+09:00")

	if err != nil {
		t.Errorf("Fail to parse err: %v", err)
	}

	post := &Post{
		ID:        1,
		Title:     "メモのタイトル",
		Body:      "メモの本文",
		Draft:     false,
		Archived:  false,
		URL:       "https://kray.docbase.io/posts/1",
		CreatedAt: ti,
		Tags: []Tag{
			Tag{Name: "rails"},
			Tag{Name: "ruby"},
		},
		Scope:      "group",
		SharingURL: "https://docbase.io/posts/1/sharing/abcdefgh-0e81-4567-9876-1234567890ab",
		User: SimpleUser{
			ID:              1,
			Name:            "danny",
			ProfileImageURL: "https://image.docbase.io/uploads/aaa.gif",
		},
		StarsCount:    1,
		GoodJobsCount: 2,
		Comments:      []Comment{},
		Groups: []SimpleGroup{
			SimpleGroup{
				ID:   1,
				Name: "DocBase",
			},
		},
		Attachments: []Attachment{Attachment{
			ID:        "461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv",
			Name:      "uploadfile.csv",
			Size:      18786,
			URL:       "https://kray.docbase.io/file_attachments/461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv",
			Markdown:  "[![csv](/images/file_icons/csv.svg) uploadfile.jpg](https://kray.docbase.io/uploads/461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv)",
			CreatedAt: ti,
		},
		},
	}

	mux.HandleFunc(fmt.Sprintf("/posts/%d", post.ID), func(w http.ResponseWriter, r *http.Request) {
		u, _ := url.Parse(fmt.Sprintf("/posts/%d", post.ID))
		want := u.String()
		// TODO impl test function
		if got := r.URL.String(); got != want {
			t.Errorf("URL: got %v, want %v", got, want)
		}
		testMethod(t, r, "GET")

		fmt.Fprint(w, testutil.LoadFixture(t, "post-detail-response.json"))
	})

	getRes, _, err := client.Posts.Get(post.ID)

	if err != nil {
		t.Errorf("Get returned an error: %v", err)
	}

	if !reflect.DeepEqual(getRes, post) {
		t.Errorf("Get returned %+v, want %+v", getRes, post)
	}
}

func TestPostService_Update(t *testing.T) {
	setup()
	defer teardown()

	ti, err := time.Parse(time.RFC3339, "2020-03-27T09:25:09+09:00")

	if err != nil {
		t.Errorf("Fail to parse err: %v", err)
	}

	post := &Post{
		ID:        1,
		Title:     "メモのタイトル",
		Body:      "メモの本文",
		Draft:     false,
		Archived:  false,
		URL:       "https://kray.docbase.io/posts/1",
		CreatedAt: ti,
		Tags: []Tag{
			Tag{Name: "rails"},
			Tag{Name: "ruby"},
		},
		Scope:      "group",
		SharingURL: "https://docbase.io/posts/1/sharing/abcdefgh-0e81-4567-9876-1234567890ab",
		User: SimpleUser{
			ID:              1,
			Name:            "danny",
			ProfileImageURL: "https://image.docbase.io/uploads/aaa.gif",
		},
		StarsCount:    1,
		GoodJobsCount: 2,
		Comments:      []Comment{},
		Groups: []SimpleGroup{
			SimpleGroup{
				ID:   1,
				Name: "DocBase",
			},
		},
		Attachments: []Attachment{Attachment{
			ID:        "461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv",
			Name:      "uploadfile.csv",
			Size:      18786,
			URL:       "https://kray.docbase.io/file_attachments/461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv",
			Markdown:  "[![csv](/images/file_icons/csv.svg) uploadfile.jpg](https://kray.docbase.io/uploads/461d38b9-8c22-4222-a6a2-a6f2ce98ec3a.csv)",
			CreatedAt: ti,
		},
		},
	}

	mux.HandleFunc(fmt.Sprintf("/posts/%d", post.ID), func(w http.ResponseWriter, r *http.Request) {
		u, _ := url.Parse(fmt.Sprintf("/posts/%d", post.ID))

		want := u.String()
		if got := r.URL.String(); got != want {
			t.Errorf("URL: got %v, want %v", got, want)
		}

		testMethod(t, r, "PATCH")
		fmt.Fprint(w, testutil.LoadFixture(t, "post-detail-response.json"))
	})

	// TDOO どんなリクエストボディでも固定レスポンス返してしまうので、検証はさみたい
	mReq := &PostUpdateRequest{}

	res, _, err := client.Posts.Update(post.ID, mReq)

	if err != nil {
		t.Errorf("Get returned an error: %v", err)
	}

	if !reflect.DeepEqual(res, post) {
		t.Errorf("Get returned %+v, want %+v", res, post)
	}
}

func TestPostService_Archive(t *testing.T) {
	setup()
	defer teardown()

	post := &Post{
		ID: 1,
	}

	mux.HandleFunc(fmt.Sprintf("/posts/%d/archive", post.ID), func(w http.ResponseWriter, r *http.Request) {
		u, _ := url.Parse(fmt.Sprintf("/posts/%d/archive", post.ID))

		want := u.String()
		if got := r.URL.String(); got != want {
			t.Errorf("URL: got %v, want %v", got, want)
		}
		testMethod(t, r, "PUT")

		fmt.Fprint(w, `{}`)
	})

	resp, err := client.Posts.Archive(post.ID)

	if err != nil {
		t.Errorf("Archive returned an error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Download.Create request code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}
}

func TestPostService_Unarchive(t *testing.T) {
	setup()
	defer teardown()

	post := &Post{
		ID: 1,
	}

	mux.HandleFunc(fmt.Sprintf("/posts/%d/unarchive", post.ID), func(w http.ResponseWriter, r *http.Request) {
		u, _ := url.Parse(fmt.Sprintf("/posts/%d/unarchive", post.ID))

		want := u.String()
		if got := r.URL.String(); got != want {
			t.Errorf("URL: got %v, want %v", got, want)
		}

		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{}`)
	})

	resp, err := client.Posts.Unarchive(post.ID)

	if err != nil {
		t.Errorf("Unarchive returned an error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Post Unarchive request code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}
}
