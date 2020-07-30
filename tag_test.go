package docbase

import (
	"fmt"
	"github.com/hayashiki/docbase-go/testutil"
	"net/http"
	"reflect"
	"testing"
)

func TestTagService_List(t *testing.T) {
	setup()
	defer teardown()

	tagSvc := NewTagService(client)

	mux.HandleFunc("/tags", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testutil.LoadFixture(t, "tag-list-response.json"))
	})

	tag1 := Tag{
		Name: "ruby",
	}

	tag2 := Tag{
		Name: "rails",
	}

	tags, _, _ := tagSvc.List()

	want := &TagListResponse{tag1, tag2}
	if !reflect.DeepEqual(tags, want) {
		t.Errorf("Users returned %+v, want %+v", tags, want)
	}
}
