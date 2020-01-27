package tpb

import (
	"context"
	"strings"
)

// Client represent a Client used to make Search
type Client struct {
	Endpoint string
}

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

// New return a new Client
func New(endpoint string) *Client {
	return &Client{
		Endpoint: strings.TrimRight(endpoint, "/"),
	}
}

// Search will search Torrents
func (c *Client) Search(ctx context.Context, search string, opts *Options) ([]*Torrent, error) {
	if opts == nil {
		opts = &DefaultOptions
	}
	url := c.Endpoint + "/search/" + search + opts.String()
	return c.fetchTorrents(ctx, url)
}

// User lists the torrents uploaded by a user
func (c *Client) User(ctx context.Context, user string, opts *Options) ([]*Torrent, error) {
	if opts == nil {
		opts = &DefaultOptions
	}
	url := c.Endpoint + "/user/" + user + opts.String()
	return c.fetchTorrents(ctx, url)
}
