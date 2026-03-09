package tpb

import (
	"testing"
	"time"
)

func TestMagnet(t *testing.T) {
	torrent := &Torrent{
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
	expectedMagnet := `magnet:?xt=urn:btih:363BC69191230430C6758318D196CCD61DB61B647&dn=Big Buck Bunny&tr=udp://tracker.opentrackr.org:1337/announce&tr=udp://open.tracker.cl:1337/announce&tr=udp://tracker.openbittorrent.com:6969/announce&tr=udp://opentracker.i2p.rocks:6969/announce&tr=udp://tracker.torrent.eu.org:451/announce&tr=udp://open.stealth.si:80/announce`
	magnet := torrent.Magnet()
	if magnet != expectedMagnet {
		t.Fatalf("expected magnet %q, got %q", expectedMagnet, magnet)
	}
}

func TestMagnetFromAPI(t *testing.T) {
	apiMagnet := "magnet:?xt=urn:btih:AABBCC&dn=Some+Torrent&tr=udp://tracker.example.com:6969/announce"
	torrent := &Torrent{
		InfoHash:   "AABBCC",
		Name:       "Some Torrent",
		MagnetLink: apiMagnet,
	}
	if got := torrent.Magnet(); got != apiMagnet {
		t.Fatalf("expected API magnet %q, got %q", apiMagnet, got)
	}
}
