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
	expectedMagnet := `magnet:?xt=urn:btih:363BC69191230430C6758318D196CCD61DB61B647&dn=Big Buck Bunny&tr=udp://tracker.coppersurfer.tk:6969/announce&tr=udp://9.rarbg.to:2920/announce&tr=udp://tracker.opentrackr.org:1337&tr=udp://tracker.internetwarriors.net:1337/announce&tr=udp://tracker.leechers-paradise.org:6969/announce&tr=udp://tracker.coppersurfer.tk:6969/announce&tr=udp://tracker.pirateparty.gr:6969/announce&tr=udp://tracker.cyberia.is:6969/announce`
	magnet := torrent.Magnet()
	if magnet != expectedMagnet {
		t.Fatalf("expected magnet %q, got %q", expectedMagnet, magnet)
	}
}
