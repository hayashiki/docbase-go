package main

import (
	"fmt"
	"github.com/hayashiki/docbase-go"
	"log"
	"os"
)

func main() {
	tagList()
	//memoGet()
	//memoUpdate()
	//memoDelete()
	//memoCreate()
	//get()
	//userList()
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
	memoSvc := docbase.NewMemoService(cli)

	res, _, _ := memoSvc.Get(1465823)

	log.Printf("response is %#v", res.Title)
	log.Printf("response is %#v", res)
}

func memoUpdate() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	memoSvc := docbase.NewMemoService(cli)

	// := memoSvc.Get("1465823")

	memoReq := &docbase.MemoRequest{
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
	memoSvc := docbase.NewMemoService(cli)

	res,_ := memoSvc.Delete("1465823")

	log.Printf("response is %#v", res)
}

func memoCreate() {
	cli := docbase.NewClient(nil, os.Getenv("DOCBASE_TEAM"), os.Getenv("DOCBASE_TOKEN"))
	memoSvc := docbase.NewMemoService(cli)

	memoReq := &docbase.MemoRequest{
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
