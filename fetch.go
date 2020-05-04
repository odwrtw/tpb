package tpb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) fetchTorrents(ctx context.Context, path string) ([]*Torrent, error) {
	var torrents []*Torrent
	err := c.fetch(ctx, path, &torrents)
	if err != nil {
		return nil, err
	}
	// If the search returned no results, there will be one "fake" torrent
	// in the results which name is: "No results returned"
	if len(torrents) == 1 && torrents[0].Name == "No results returned" {
		return []*Torrent{}, nil
	}

	return torrents, nil
}

// fetch will try to GET the path, trying all the endpoints if needed, and
// unmarshal the results in the data interface
func (c *Client) fetch(ctx context.Context, path string, data interface{}) error {
	var err error
	for i := 0; i < c.MaxTries; i++ {
		endpoint := c.endpoints.best()
		if endpoint == nil {
			return ErrMissingEndpoint
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, c.EndpointTimeout)
		defer cancel()

		err = get(timeoutCtx, endpoint.baseURL+path, &data)
		if err == nil {
			return nil
		}
		endpoint.lastFailure = time.Now()

		// Stop if the global context is done
		if ctx.Err() != nil {
			err = ctx.Err()
			break
		}
	}

	return err
}

// get will GET the url and unmarshal the results in the data interface
func get(ctx context.Context, url string, data interface{}) error {
	var err error
	done := make(chan struct{})

	go func() {
		defer close(done)
		var resp *http.Response
		resp, err = http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("got status %d when making the request", resp.StatusCode)
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&data)
	}()

	select {
	case <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
