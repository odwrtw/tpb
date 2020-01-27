package tpb

import (
	"context"
	"fmt"
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

// SearchOptions represent the options for a Search lookup
type SearchOptions struct {
	Key      string
	OrderBy  OrderBy
	Page     int
	Sort     SortOrder
	Category TorrentCategory
}

func (so *SearchOptions) String() string {
	return fmt.Sprintf(
		"/search/%s/%d/%s/%d/",
		so.Key,
		so.Page,
		mapOrderBy(so.OrderBy, so.Sort),
		int(so.Category),
	)
}

// UserOptions represent the options for a User lookup
type UserOptions struct {
	OrderBy OrderBy
	Page    int
	Sort    SortOrder
	User    string
}

func (uo *UserOptions) String() string {
	return fmt.Sprintf(
		"/user/%s/%d/%s/0/",
		uo.User,
		uo.Page,
		mapOrderBy(uo.OrderBy, uo.Sort),
	)
}

// SearchWithCtx will search Torrents with a context
func (c *Client) SearchWithCtx(ctx context.Context, opt SearchOptions) ([]*Torrent, error) {
	url := c.Endpoint + opt.String()
	return c.fetchTorrents(ctx, url)
}

// Search will search Torrents
func (c *Client) Search(opt SearchOptions) ([]*Torrent, error) {
	ctx := context.Background()
	return c.SearchWithCtx(ctx, opt)
}

// UserWithContext will return a user's Torrents with a context
func (c *Client) UserWithContext(ctx context.Context, opt UserOptions) ([]*Torrent, error) {
	url := c.Endpoint + opt.String()
	return c.fetchTorrents(ctx, url)
}

// User will return a user's Torrents
func (c *Client) User(opt UserOptions) ([]*Torrent, error) {
	ctx := context.Background()
	return c.UserWithContext(ctx, opt)
}

// FilterByUsers will filter the results and return only those from the given
// users list
func FilterByUsers(torrents []*Torrent, users []string) []*Torrent {
	filteredTorrents := []*Torrent{}
	for _, t := range torrents {
		for _, u := range users {
			if t.User == u {
				filteredTorrents = append(filteredTorrents, t)
			}
		}
	}
	return filteredTorrents
}
