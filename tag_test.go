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

	mux.HandleFunc("/tags", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testutil.LoadFixture(t, "tag-list-response.json"))
	})

	tag1 := Tag{
		Name: "ruby",
	}

	tag2 := Tag{
		Name: "rails",
	}

	tags, resp, err := client.Tags.List()

	if err != nil {
		t.Errorf("Shouldn't have returned an error: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Comment Create request code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}

	want := &TagListResponse{tag1, tag2}
	if !reflect.DeepEqual(tags, want) {
		t.Errorf("Users returned %+v, want %+v", tags, want)
	}
}
