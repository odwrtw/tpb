package tpb

import (
	"reflect"
	"testing"
	"time"
)

func TestNewEndpoints(t *testing.T) {
	urls := []string{
		"http://test1.com/",
		"http://test2.com",
		"http://test3.com/",
		"http://test4.com",
	}

	got := newEndpoints(urls)
	expected := endpoints{
		{baseURL: urls[0][:16]},
		{baseURL: urls[1]},
		{baseURL: urls[2][:16]},
		{baseURL: urls[3]},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected:\n%v\ngot:\n%v", expected, got)
	}
}

func TestBestEndpoint(t *testing.T) {
	e1 := &endpoint{lastFailure: time.Now().Add(-1 * time.Minute)}
	e2 := &endpoint{lastFailure: time.Now().Add(-2 * time.Minute)}
	e3 := &endpoint{}
	e4 := &endpoint{lastFailure: time.Now()}

	input := endpoints{e1, e2, e3, e4}
	best := input.best()

	if best != e3 {
		t.Fatalf("failed to find the best endpoint")
	}
}

func TestNoEndpoint(t *testing.T) {
	e := endpoints{}
	if e.best() != nil {
		t.Fatalf("no endpoints should mean no best endpoint")
	}
}
