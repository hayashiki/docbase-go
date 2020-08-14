package docbase

import (
	"fmt"
	"github.com/hayashiki/docbase-go/testutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGroupService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testutil.LoadFixture(t, "group-list-response.json"))
	})

	group1 := SimpleGroup{
		ID:   1,
		Name: "DocBase",
	}

	group2 := SimpleGroup{
		ID:   2,
		Name: "kray-internal",
	}

	opts := &GroupListOptions{
		PerPage: 5,
		Page:    1,
		Name:    "query",
	}

	groups, resp, err := client.Groups.List(opts)

	want := &GroupListResponse{group1, group2}

	if err != nil {
		t.Errorf("Shouldn't have returned an error: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Comment Create request code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}

	if !reflect.DeepEqual(groups, want) {
		t.Errorf("GroupList returned %+v, want %+v", groups, want)
	}
}

func TestGroupService_Get(t *testing.T) {
	setup()
	defer teardown()

	ti, err := time.Parse(time.RFC3339, "2020-03-27T09:25:09+09:00")

	if err != nil {
		t.Errorf("Fail to parse err: %v", err)
	}

	group := &Group{
		ID:             1,
		Name:           "グループ1",
		Description:    "APIで作ったグループ",
		PostsCount:     0,
		LastActivityAt: ti,
		CreatedAt:      ti,
		Users: []SimpleUser{
			SimpleUser{
				ID:              1,
				Name:            "docbaseman",
				ProfileImageURL: "https://image.docbase.io/uploads/aaa.gif",
			},
		},
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%d", group.ID), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testutil.LoadFixture(t, "group-response.json"))
	})

	actual, _, _ := client.Groups.Get(group.ID)

	if !reflect.DeepEqual(actual, group) {
		t.Errorf("Group returned %+v, want %+v", actual, group)
	}
}

func TestGroupCli_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &GroupCreateRequest{
		Name:        "DocBase",
		Description: "DocBase",
	}

	mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		fmt.Fprint(w, testutil.LoadFixture(t, "group-response.json"))
	})

	group, resp, err := client.Groups.Create(createRequest)

	if err != nil {
		t.Errorf("Fail to create group request err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Create group response code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}

	ti, err := time.Parse(time.RFC3339, "2020-03-27T09:25:09+09:00")

	want := &Group{
		ID:             1,
		Name:           "グループ1",
		Description:    "APIで作ったグループ",
		PostsCount:     0,
		LastActivityAt: ti,
		CreatedAt:      ti,
		Users:          []SimpleUser{
			SimpleUser{
				ID:              1,
				Name:            "docbaseman",
				ProfileImageURL: "https://image.docbase.io/uploads/aaa.gif",
			},
		},
	}

	// TODO
	if !reflect.DeepEqual(want, group) {
		t.Errorf("Request response: %+v, want %+v", group, want)
	}
}
