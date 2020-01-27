package tpb

import (
	"context"

	"github.com/gocolly/colly/v2"
)

func (c *Client) fetchTorrents(ctx context.Context, url string) ([]*Torrent, error) {
	var err error
	torrents := []*Torrent{}
	done := make(chan struct{})

	co := colly.NewCollector()
	co.OnHTML("#searchResult", func(e *colly.HTMLElement) {
		e.ForEach("tbody > tr", func(i int, se *colly.HTMLElement) {
			if err != nil {
				// If an error occured already, don't do any more work
				return
			}

			data := rawData{}
			se.Unmarshal(&data)

			torrent, err := data.parse()
			if err != nil {
				return
			}

			torrents = append(torrents, torrent)
		})
		done <- struct{}{}
	})

	go func() {
		err = co.Visit(url)
		close(done)
	}()

	select {
	case <-done:
		return torrents, err
	case <-ctx.Done():
		return torrents, ctx.Err()
	}
}
