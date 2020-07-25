package docbase

import (
	"fmt"
	"github.com/hayashiki/docbase-go/testutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestUserService_List(t *testing.T) {
	setup()
	defer teardown()

	userSvc := NewUserService(client)

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testutil.LoadFixture(t, "user-list-response.json"))
	})

	ti, err := time.Parse(time.RFC3339, "2020-03-27T09:25:09+09:00")

	if err != nil {
		t.Errorf("Fail to parse err: %v", err)
	}

	user1 := &User{
		ID:                    1,
		Name:                  "ドックベースマン",
		Username:              "docbaseman",
		ProfileImageURL:       "https://image.docbase.io/uploads/aaa.gif",
		Role:                  "owner",
		PostsCount:            2,
		LastAccessTime:        ti, // 2019-02-18T11:52:56.000+09:00
		TwoStepAuthentication: false,
		Groups:                []UGroup{
			UGroup{
				ID: 1,
				Name: "グループ1",
			},
		},
	}

	user2 := &User{
		ID:                    2,
		Name:                  "ドックベースウーマン",
		Username:              "docbasewoman",
		ProfileImageURL:       "https://image.docbase.io/uploads/aaa.gif",
		Role:                  "admin",
		PostsCount:            3,
		LastAccessTime:        ti, //2019-02-18T11:52:56.000+09:00
		TwoStepAuthentication: false,
		Groups:                []UGroup{},
	}

	user3 := &User{
		ID:                    3,
		Name:                  "ドックべーサー",
		Username:              "docbaser",
		ProfileImageURL:       "https://image.docbase.io/uploads/aaa.gif",
		Role:                  "user",
		PostsCount:            5,
		LastAccessTime:        ti, //"2019-02-18T11:52:56.000+09:00"
		TwoStepAuthentication: false,
		Groups:                []UGroup{},
	}

	opts := &UserListOptions{
		PerPage: 5,
		Page: 1,
		Q: "query",
	}

	users, _, _ := userSvc.List(opts)

	want := &[]User{*user1, *user2, *user3}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("Users returned %+v, want %+v", users, want)
	}
}
