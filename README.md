# Golang scrapper for thepiratebay

[![Build Status](https://travis-ci.org/odwrtw/tpb.svg?branch=master)](https://travis-ci.org/odwrtw/tpb)
[![Go Report Card](https://goreportcard.com/badge/github.com/odwrtw/tpb)](https://goreportcard.com/report/github.com/odwrtw/tpb)
[![GoDoc](https://godoc.org/github.com/odwrtw/tpb?status.png)](http://godoc.org/github.com/odwrtw/tpb)
[![Coverage Status](https://coveralls.io/repos/github/odwrtw/tpb/badge.svg?branch=master)](https://coveralls.io/github/odwrtw/tpb?branch=master)

## Exemple

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/odwrtw/tpb"
)

func main() {
    if err := run(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func run() error {
    // The client supports multiple endpoints and will try to use one that
    // works
    client := tpb.NewWithEndpoints([]string{
        "https://thepiratebay.org",
        "https://mypirate.cc",
        "https://piratebay.life",
    })

    // You can create a context to cancel the search
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    // You can add search options or nil for the default options
    torrents, err := client.Search(ctx, "Ubuntu", nil)
    if err != nil {
        return err
    }

    for _, t := range torrents {
        fmt.Println("--------------")
        fmt.Printf("%s\nUploaded by %q (%d seeders / %d leechers)\n",
            t.Name,
            t.User,
            t.Seeders,
            t.Leechers,
        )
    }

    return nil
}
```
