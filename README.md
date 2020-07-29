
![Go](https://github.com/hayashiki/docbase-go/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/hayashiki/docbase-go/branch/develop/graph/badge.svg)](https://codecov.io/gh/hayashiki/docbase-go)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

# docbase-go

docbase-go is unofficial Go wrapper tool for Docbase API.

# Installation

```
go get 
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
  fmt.Printf("%+v", posts)
}

```

# API

## Posts

``` go
// Get the information about the post list
posts, resp, err := client.Posts.List(&docbase.PostListOptions{})
```

# Note

[Here is the full API.](https://help.docbase.io/posts/45703)

# License

