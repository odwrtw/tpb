package tpb

import (
	"context"
	"errors"
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
	// DefaultTimeout defines the timeout to use by endpoint
	DefaultTimeout time.Duration
}

// New return a new client with a the given endpoint(s)
func New(endpoints ...string) *Client {
	return &Client{
		endpoints:      newEndpoints(endpoints...),
		MaxTries:       len(endpoints),
		DefaultTimeout: defaultTimeout,
	}
}

func (c *Client) fetch(ctx context.Context, path string) ([]*Torrent, error) {
	var err error
	for i := 0; i < c.MaxTries; i++ {
		endpoint := c.endpoints.best()
		if endpoint == nil {
			return nil, ErrMissingEndpoint
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, c.DefaultTimeout)
		defer cancel()

		var torrents []*Torrent
		torrents, err = fetchTorrents(timeoutCtx, endpoint.baseURL+path)
		if err == nil {
			return torrents, nil
		}
		endpoint.lastFailure = time.Now()

		// Stop if the error is a parsing error
		if errors.As(err, &ParserError{}) {
			break
		}

		// Stop if the global context is done
		if ctx.Err() != nil {
			err = ctx.Err()
			break
		}
	}

	return nil, err
}

// Search will search Torrents
func (c *Client) Search(ctx context.Context, search string, opts *Options) ([]*Torrent, error) {
	if opts == nil {
		opts = &DefaultOptions
	}
	path := "/search/" + search + opts.String()
	return c.fetch(ctx, path)
}

// User lists the torrents uploaded by a user
func (c *Client) User(ctx context.Context, user string, opts *Options) ([]*Torrent, error) {
	if opts == nil {
		opts = &DefaultOptions
	}
	path := "/user/" + user + opts.String()
	return c.fetch(ctx, path)
}
