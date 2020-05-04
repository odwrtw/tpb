package tpb

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

var defaultTimeout = 10 * time.Second

// ErrMissingEndpoint is the error returned if there's no endpoint
var ErrMissingEndpoint = errors.New("tpb: missing endpoint")

// Client represent a Client used to make Search
type Client struct {
	endpoints endpoints
	// MaxTries represent the number total number of endpoint to try to find a
	// result
	MaxTries int
	// EndpointTimeout defines the timeout to use by endpoint
	EndpointTimeout time.Duration
}

// New return a new client with a the given endpoint(s)
func New(endpoints ...string) *Client {
	return &Client{
		endpoints:       newEndpoints(endpoints...),
		MaxTries:        len(endpoints),
		EndpointTimeout: defaultTimeout,
	}
}

// Search will search torrents
func (c *Client) Search(ctx context.Context, search string, opts *SearchOptions) ([]*Torrent, error) {
	if opts == nil {
		opts = &SearchOptions{
			Category: All,
		}
	}
	v := opts.Category.URLValue()
	v.Add("q", search)
	path := "/q.php?" + v.Encode()
	return c.fetchTorrents(ctx, path)
}

// User lists the torrents uploaded by a user
func (c *Client) User(ctx context.Context, user string, opts *UserOptions) ([]*Torrent, error) {
	if opts == nil {
		opts = &UserOptions{
			Page: 0,
		}
	}
	query := fmt.Sprintf("user:%s:%d", user, opts.Page)
	v := url.Values{}
	v.Add("q", query)
	path := "/q.php?" + v.Encode()
	return c.fetchTorrents(ctx, path)
}

// Category lists the latest torrents for a category
func (c *Client) Category(ctx context.Context, category TorrentCategory, opts *CategoryOptions) ([]*Torrent, error) {
	if opts == nil {
		opts = &CategoryOptions{
			Page: 0,
		}
	}
	query := fmt.Sprintf("category:%d:%d", category, opts.Page)
	v := url.Values{}
	v.Add("q", query)
	path := "/q.php?" + v.Encode()
	return c.fetchTorrents(ctx, path)
}

// Top100 lists the Top100 torrents
func (c *Client) Top100(ctx context.Context, cat TorrentCategory) ([]*Torrent, error) {
	var catString = "all"
	if cat != All {
		catString = strconv.Itoa(int(cat))
	}
	path := "/precompiled/data_top100_" + catString + ".json"
	return c.fetchTorrents(ctx, path)
}

// Infos gets infos about a given torrent ID
func (c *Client) Infos(ctx context.Context, id int) (*Torrent, error) {
	t := &Torrent{}
	v := url.Values{}
	v.Add("id", strconv.Itoa(id))
	path := "/t.php?" + v.Encode()
	return t, c.fetch(ctx, path, &t)
}

// TODO: Implement FileList
// curl https://apibay.org/f.php?id=36120091
