package docbase

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGroupUserCli_Create(t *testing.T) {
	setup()
	defer teardown()

	req := &GroupUserCreateRequest{
		UserIDs: []int{43492},
	}

	group := &SimpleGroup{
		ID:   1,
		Name: "DocBase",
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%d/users", group.ID), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	})

	resp, err := client.GroupUsers.Create(group.ID, req)

	if err != nil {
		t.Errorf("Fail to request err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("UserGroup Delete request code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}
}

func TestGroupUserCli_Delete(t *testing.T) {
	setup()
	defer teardown()

	req := &GroupUserCreateRequest{
		UserIDs: []int{43492},
	}

	group := &SimpleGroup{
		ID:   1,
		Name: "DocBase",
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%d/users", group.ID), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	})

	resp, err := client.GroupUsers.Delete(group.ID, req)

	if err != nil {
		t.Errorf("Fail to request err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("UserGroup Delete request code = %v, expected %v", resp.StatusCode, http.StatusOK)
	}
}
