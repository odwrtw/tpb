package tpb

// Torrent represents a Torrent
type Torrent struct {
	Name        string
	Peers       int
	Seeds       int
	User        string
	Magnet      string
	Size        uint64
	Category    TorrentCategory
	SubCategory TorrentCategory
}
