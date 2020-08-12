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

	groupSvc := NewGroupService(client)

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

	groups, _, _ := groupSvc.List(opts)

	want := &GroupListResponse{group1, group2}
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

	groupSvc := NewGroupService(client)

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

	actual, _, _ := groupSvc.Get(group.ID)

	if !reflect.DeepEqual(actual, group) {
		t.Errorf("Group returned %+v, want %+v", actual, group)
	}
}

func TestGroupCli_Create(t *testing.T) {
	setup()
	defer teardown()

	groupSvc := NewGroupService(client)

	createReqest := &GroupCreateRequest{
		Name: "DocBase",
		Description:   "DocBase",
	}

	mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	})

	_, resp, err := groupSvc.Create(createReqest)

	if err != nil {
		t.Errorf("Fail to create group request err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Create group response code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}

}

func TestGroupService_AddUser(t *testing.T) {
	setup()
	defer teardown()

	groupSvc := NewGroupService(client)

	req := &GroupAddRequest{
		UserIDs: []int{43492},
	}

	group := &SimpleGroup{
		ID:   1,
		Name: "DocBase",
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%d/users", group.ID), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	})

	_, err := groupSvc.AddUser(group.ID, req)

	if err != nil {
		t.Errorf("Fail to request err: %v", err)
	}
}

func TestGroupService_RemoveUser(t *testing.T) {
	setup()
	defer teardown()

	groupSvc := NewGroupService(client)

	req := &GroupAddRequest{
		UserIDs: []int{43492},
	}

	group := &SimpleGroup{
		ID:   1,
		Name: "DocBase",
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%d/users", group.ID), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	})

	_, err := groupSvc.RemoveUser(group.ID, req)

	if err != nil {
		t.Errorf("Fail to request err: %v", err)
	}
}
