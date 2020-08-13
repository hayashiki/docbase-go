package docbase

import (
	"fmt"
	"github.com/hayashiki/docbase-go/testutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestCommentService_Create(t *testing.T) {
	setup()
	defer teardown()

	post := &Post{
		ID: 1,
	}

	mux.HandleFunc(fmt.Sprintf("/posts/%d/comments", post.ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		fmt.Fprint(w, testutil.LoadFixture(t, "comment-response.json"))
	})

	cReq := &CommentCreateRequest{}

	ti, err := time.Parse(time.RFC3339, "2020-03-27T09:25:09+09:00")

	want := &Comment{
		ID:   1,
		Body: "コメント",
		CreatedAt: ti,
		SimpleUser: SimpleUser{
			ID: 1,
			Name: "danny",
			ProfileImageURL: "https://image.docbase.io/uploads/aaa.gif",
		},
	}

	comment, resp, err := client.Comments.Create(post.ID, cReq)

	if err != nil {
		t.Errorf("Shouldn't have returned an error: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Comment Create request code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}

	if !reflect.DeepEqual(want, comment) {
		t.Errorf("Request response: %+v, want %+v", comment, want)
	}
}

func TestCommentCli_Create_Error(t *testing.T) {
	setup()
	defer teardown()

	post := &Post{
		ID: 1,
	}

	mux.HandleFunc(fmt.Sprintf("/posts/%d/comments", post.ID), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{
		  "errors":
		  [
				"Bad request"
		  ]
		}`)
	})

	cReq := &CommentCreateRequest{}

	_, resp, err := client.Comments.Create(post.ID, cReq)

	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Errorf("Error should be of type ErrorResponse but is %v: %+v", reflect.TypeOf(err), err)
	}

	want := "Bad Request"
	if got := errResp.Messages[0]; want != got {
		t.Errorf("Error message: %v, want %v", got, want)
	}

	if got, want := resp.StatusCode, http.StatusBadRequest; want != got {
		t.Errorf("Status code: %d, want %d", got, want)
	}
}

func TestCommentService_Delete(t *testing.T) {
	setup()
	defer teardown()

	comment := &Comment{
		ID:   1,
	}

	mux.HandleFunc(fmt.Sprintf("/comments/%d", comment.ID), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		testMethod(t, r, "DELETE")

		fmt.Fprint(w, `{}`)
	})

	resp, err := client.Comments.Delete(comment.ID)

	if err != nil {
		//t.Errorf("Shouldn't have returned an error: %+v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Comment Delete request code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}
}
