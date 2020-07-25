package docbase

import (
	"fmt"
	"github.com/hayashiki/docbase-go/testutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestCommentService_Create(t *testing.T) {
	setup()
	defer teardown()

	post := &Memo{
		ID: 1,
	}

	mux.HandleFunc(fmt.Sprintf("/posts/%d/comments", post.ID), func(w http.ResponseWriter, r *http.Request) {
		log.Printf("memo")
		fmt.Fprint(w, testutil.LoadFixture(t, "comment-response.json"))
	})

	commentSvc := NewCommentService(client)

	cReq := &CommentRequest{}

	want := &CommentResponse{
		Body:   "コメント",
	}

	comment, _, err := commentSvc.Create(strconv.Itoa(post.ID), cReq)

	if err != nil {
		t.Errorf("Shouldn't have returned an error: %+v", err)
	}

	if !reflect.DeepEqual(want, comment) {
		t.Errorf("Request response: %+v, want %+v", comment, want)
	}
}

func TestCommentService_Delete(t *testing.T) {

}
