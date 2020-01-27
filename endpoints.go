package tpb

import (
	"sort"
	"strings"
	"time"
)

type endpoint struct {
	baseURL     string
	lastFailure time.Time
}

type endpoints []*endpoint

func newEndpoints(urls []string) endpoints {
	e := make([]*endpoint, len(urls))
	for i, u := range urls {
		e[i] = &endpoint{
			baseURL: strings.TrimRight(u, "/"),
		}
	}

	return e
}

func (e endpoints) best() *endpoint {
	if len(e) == 0 {
		return nil
	}

	sort.SliceStable(e, func(i, j int) bool {
		return e[i].lastFailure.Before(e[j].lastFailure)
	})
	return e[0]
}
