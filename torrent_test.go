package tpb

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var rawSearchByUser = `
{
  "id": "18837600",
  "name": "ArchLinux",
  "info_hash": "B137DE1DF926E787FE263D4187B34B23",
  "leechers": "39",
  "seeders": "96",
  "num_files": "2",
  "size": "45486511610",
  "username": "user1",
  "added": "1509051120",
  "status": "vip",
  "category": "303",
  "imdb": ""
}`
var rawSearch = `
{
   "id": "6665688",
   "name": "Big Buck Bunny",
   "info_hash": "363BC69191230430C6758318D196CCD61DB61B647",
   "leechers": "117",
   "seeders": "718",
   "num_files": "5",
   "size": "738095006",
   "username": "bob",
   "added": "1509051120",
   "status": "member",
   "category": "201",
   "imdb": "tt1254207"
}`
var rawTorrentInfo = `
{
   "id": 6665688,
   "category": 201,
   "status": "member",
   "name": "Big Buck Bunny",
   "num_files": 5,
   "size": 738095006,
   "seeders": 718,
   "leechers": 117,
   "username": "bob",
   "added": 1509051120,
   "descr": "description of Big Buck Bunny",
   "imdb": "tt1254207",
   "language": 1,
   "textlanguage": 1,
   "info_hash": "363BC69191230430C6758318D196CCD61DB61B647"
}`
var rawTorrentPrecompiled = `
{
   "id": 6665688,
   "info_hash": "363BC69191230430C6758318D196CCD61DB61B647",
   "category": 201,
   "name": "Big Buck Bunny",
   "status": "member",
   "num_files": 5,
   "size": 738095006,
   "seeders": 718,
   "leechers": 117,
   "username": "bob",
   "added": 1509051120,
   "anon": 0,
   "imdb": null
}`

func TestSearch(t *testing.T) {
	var requestURI string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI = r.RequestURI
		fmt.Fprintf(w, "%s", fmt.Sprintf("[%s]", rawSearch))
	}))
	defer ts.Close()

	expected := []*Torrent{
		&Torrent{
			ID:       6665688,
			Name:     "Big Buck Bunny",
			InfoHash: "363BC69191230430C6758318D196CCD61DB61B647",
			Leechers: 117,
			Seeders:  718,
			NumFiles: 5,
			Size:     738095006,
			User:     "bob",
			Status:   Member,
			Category: VideoMovies,
			ImdbID:   "tt1254207",
			Added:    time.Unix(1509051120, 0),
		},
	}

	client := New(ts.URL)

	got, err := client.Search(context.Background(), "Big Buck Bunny", nil)
	if err != nil {
		t.Fatalf("got error while searching %q", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected: \n%+v\n, got \n%+v", expected, got)
	}

	expectedRequestURI := "/q.php?cat=0&q=Big+Buck+Bunny"
	if requestURI != expectedRequestURI {
		t.Fatalf("expected URL %q, got %q", expectedRequestURI, requestURI)
	}
}

func TestUser(t *testing.T) {
	var requestURI string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI = r.RequestURI
		fmt.Fprintf(w, "%s", fmt.Sprintf("[%s]", rawSearchByUser))
	}))
	defer ts.Close()

	expected := []*Torrent{
		&Torrent{
			ID:       18837600,
			Name:     "ArchLinux",
			InfoHash: "B137DE1DF926E787FE263D4187B34B23",
			Leechers: 39,
			Seeders:  96,
			NumFiles: 2,
			Size:     45486511610,
			User:     "user1",
			Status:   VIP,
			Category: ApplicationsUNIX,
			ImdbID:   "",
			Added:    time.Unix(1509051120, 0),
		},
	}

	client := New(ts.URL)

	got, err := client.User(context.Background(), "user1", nil)
	if err != nil {
		t.Fatalf("got error while listing user's torrents %q", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected: \n%+v\n, got \n%+v", expected, got)
	}

	expectedRequestURI := "/q.php?q=" + url.QueryEscape("user:user1:0")
	if requestURI != expectedRequestURI {
		t.Fatalf("expected URL %q, got %q", expectedRequestURI, requestURI)
	}
}

func TestCategory(t *testing.T) {
	var requestURI string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI = r.RequestURI
		fmt.Fprintf(w, "%s", fmt.Sprintf("[%s]", rawSearch))
	}))
	defer ts.Close()

	expected := []*Torrent{
		&Torrent{
			ID:       6665688,
			Name:     "Big Buck Bunny",
			InfoHash: "363BC69191230430C6758318D196CCD61DB61B647",
			Leechers: 117,
			Seeders:  718,
			NumFiles: 5,
			Size:     738095006,
			User:     "bob",
			Status:   Member,
			Category: VideoMovies,
			ImdbID:   "tt1254207",
			Added:    time.Unix(1509051120, 0),
		},
	}

	client := New(ts.URL)

	got, err := client.Category(context.Background(), All, nil)
	if err != nil {
		t.Fatalf("got error while listing category's torrents %q", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected: \n%+v\n, got \n%+v", expected, got)
	}

	expectedRequestURI := "/q.php?q=" + url.QueryEscape("category:0:0")
	if requestURI != expectedRequestURI {
		t.Fatalf("expected URL %q, got %q", expectedRequestURI, requestURI)
	}
}

func TestTop100(t *testing.T) {
	var requestURI string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI = r.RequestURI
		fmt.Fprintf(w, "%s", fmt.Sprintf("[%s]", rawTorrentPrecompiled))
	}))
	defer ts.Close()

	expected := []*Torrent{
		&Torrent{
			ID:       6665688,
			Name:     "Big Buck Bunny",
			InfoHash: "363BC69191230430C6758318D196CCD61DB61B647",
			Leechers: 117,
			Seeders:  718,
			NumFiles: 5,
			Size:     738095006,
			User:     "bob",
			Status:   Member,
			Category: VideoMovies,
			ImdbID:   "",
			Added:    time.Unix(1509051120, 0),
		},
	}

	client := New(ts.URL)

	got, err := client.Top100(context.Background(), VideoMovies)
	if err != nil {
		t.Fatalf("got error while listing category's torrents %q", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected: \n%+v\n, got \n%+v", expected[0], got[0])
	}

	expectedRequestURI := "/precompiled/data_top100_201.json"
	if requestURI != expectedRequestURI {
		t.Fatalf("expected URL %q, got %q", expectedRequestURI, requestURI)
	}
}

func TestTorrentInfo(t *testing.T) {
	var requestURI string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI = r.RequestURI
		fmt.Fprintf(w, "%s", rawTorrentInfo)
	}))
	defer ts.Close()

	expected := &Torrent{
		ID:          6665688,
		Name:        "Big Buck Bunny",
		InfoHash:    "363BC69191230430C6758318D196CCD61DB61B647",
		Leechers:    117,
		Seeders:     718,
		NumFiles:    5,
		Size:        738095006,
		User:        "bob",
		Status:      Member,
		Category:    VideoMovies,
		ImdbID:      "tt1254207",
		Description: "description of Big Buck Bunny",
		Added:       time.Unix(1509051120, 0),
	}

	client := New(ts.URL)

	got, err := client.Infos(context.Background(), 6665688)
	if err != nil {
		t.Fatalf("got error while listing category's torrents %q", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected: \n%+v\n, got \n%+v", expected, got)
	}

	expectedRequestURI := "/t.php?id=6665688"
	if requestURI != expectedRequestURI {
		t.Fatalf("expected URL %q, got %q", expectedRequestURI, requestURI)
	}
}
