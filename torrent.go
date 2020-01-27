package tpb

// Torrent represents a Torrent
type Torrent struct {
	Name        string
	Leechers    int
	Seeders     int
	User        string
	Magnet      string
	Size        uint64
	Category    TorrentCategory
	SubCategory TorrentCategory
}
