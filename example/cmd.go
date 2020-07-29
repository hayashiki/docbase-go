package main

import (
	"fmt"
	"github.com/hayashiki/docbase-go"
	"log"
	"os"
)

func main() {
	//postList()
	//archive()
	//unarchive()

	//addGroupUser()
	//deleteGroupUser()
	//group()
	groupList()
	//tagList()
	//memoGet()
	//memoUpdate()
	//memoDelete()
	//memoCreate()
	//get()
	//userList()
}

func postList() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	postSvc := docbase.NewPostService(cli)
	opts := &docbase.PostListOptions{
		PerPage: 5,
		Page:    1,
		//Q:       "haya",
	}
	posts, _, _ := postSvc.List(opts)

	log.Printf("posts %#v", posts)
}

func archive() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_READTOKEN"), os.Getenv("DOCBASE_TOKEN"))
	memoSvc := docbase.NewPostService(cli)

	resp, err := memoSvc.Archive(1333773)

	log.Printf("resp %+v", resp)
	log.Printf("err %+v", err)

}

func unarchive() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	memoSvc := docbase.NewPostService(cli)

	memoSvc.Unarchive(1333773)
}

func groupList() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	groupSvc := docbase.NewGroupService(cli)

	opts := &docbase.GroupListOptions{
		Name:    "Test",
		Page:    1,
		PerPage: 2,
	}

	groups, _, _ := groupSvc.List(opts)

	log.Printf("groups %+v", groups)
}

func group() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	groupSvc := docbase.NewGroupService(cli)

	group, _, _ := groupSvc.Get(17737)

	log.Printf("groups %+v", group)
}

func addGroupUser() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	groupSvc := docbase.NewGroupService(cli)

	//ID:43492,
	req := &docbase.GroupAddRequest{
		UserIDs: []int{43492},
	}

	groupSvc.AddUser(17769, req)
	//log.Printf("groups %+v", group)

}

func deleteGroupUser() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	groupSvc := docbase.NewGroupService(cli)

	//ID:43492,
	req := &docbase.GroupAddRequest{
		UserIDs: []int{43492},
	}

	// iso 45896
	resp, err := groupSvc.RemoveUser(17769, req)
	log.Printf("groups %+v", resp)
	log.Printf("groups %+v", err)

}

func tagList() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	tagSvc := docbase.NewTagService(cli)

	tags, _, _ := tagSvc.List()

	log.Printf("tags %+v", tags)

}

func userList() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	userSvc := docbase.NewUserService(cli)
	opts := &docbase.UserListOptions{
		PerPage: 5,
		Page:    1,
		Q:       "haya",
	}
	users, _, _ := userSvc.List(opts)

	log.Printf("users %#v", users)
}

func memoGet() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	memoSvc := docbase.NewPostService(cli)

	res, _, _ := memoSvc.Get(1465823)

	log.Printf("response is %#v", res.Title)
	log.Printf("response is %#v", res)
}

func memoUpdate() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	memoSvc := docbase.NewPostService(cli)

	// := memoSvc.Get("1465823")

	memoReq := &docbase.PostRequest{
		Title: "Example memo2",
		Body:  "Example body2",
		Draft: false,
		Scope: "private",
	}

	res, _, _ := memoSvc.Update(1465823, memoReq)

	log.Printf("response is %#v", res.Title)
	log.Printf("response is %#v", res)
}

func memoDelete() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	memoSvc := docbase.NewPostService(cli)

	res, _ := memoSvc.Delete("1465823")

	log.Printf("response is %#v", res)
}

func memoCreate() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	memoSvc := docbase.NewPostService(cli)

	memoReq := &docbase.PostRequest{
		Title: "Example memo",
		Body:  "Example body",
		Draft: false,
		Scope: "private",
	}

	res, _, err := memoSvc.Create(memoReq)

	if err != nil {
		fmt.Errorf("Fail to create memo, err:%w", err.Error())
	}

	log.Printf("response is %#v", res)
}
