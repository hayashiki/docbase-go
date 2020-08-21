
![Test](https://github.com/hayashiki/docbase-go/workflows/Test/badge.svg)
[![codecov](https://codecov.io/gh/hayashiki/docbase-go/branch/develop/graph/badge.svg)](https://codecov.io/gh/hayashiki/docbase-go)
![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)

# docbase-go

Unofficial Docbase API client library, written in Go.


# Installation

```
go get github.com/hayashiki/docbase-go
```

# Example

Get all your team posts

```
package main

import (
  "fmt"
  "github.com/hayashiki/docbase-go" 
)


func main() {
  client = docbase.NewClient(nil, "your_team", "your_token")

  opts := &docbase.PostListOptions{
    PerPage: 5,
    Page:    1,
    Q:       "query something",
  }
  
  posts, resp, err := client.Posts.List(opts)
  if err != nil {
    return
  }
  fmt.Printf("resp: %+v, posts: %+v", posts)
}

```

# API

## Posts

``` go
// Get the information about the post list
posts, resp, err := client.Posts.List(&docbase.PostListOptions{})

// Get the information about the post detail
post, resp, err := client.Posts.Get(1234567)

// Create a new post
post, resp, err := client.Posts.Create(&docbase.PostRequest{})

// Update the information about the post
post, resp, err := client.Posts.Update(1234567, &docbase.PostRequest{})

// Archive the post
resp, err := client.Posts.Archive(1234567)

// Unarchive the post
resp, err := client.Posts.Unarchive(1234567)


```

## Groups

``` go

// Get the information about the group list
groups, resp, err := client.Groups.List(&docbase.GroupListOptions{})

// Get the information about the group detail
group, resp, err := client.Groups.Get(12345)

// Add the user to the group
resp, err := client.GroupUsers.Create(12345, &docbase.GroupUserCreateRequest{})

// Remove the user to the group
resp, err := client.GroupUsers.Delete(12345, &docbase.GroupUserCreateRequest{})

```

## Tags

``` go
// Get the information about the tag list
tags, resp, err := client.Tags.List()

```

## Users

``` go
// Get the information about the tag list
users, resp, err := client.Users.List(&docbase.UserListOptions{})

```

## Attachments

``` go
// Download an attachment
file, resp, err := client.Attachments.Download("8babf378-1234-5678-b62b-5a2a6c536b2b.png")

// Upload an attachment, select multiple files path
files, resp, err := client.Attachments.Upload([]string{"./testdata/test-image.jpg"})
```

# Note

[Here is the original full API.](https://help.docbase.io/posts/45703)

