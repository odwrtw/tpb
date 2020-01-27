package tpb

import (
	"context"

	"github.com/gocolly/colly/v2"
)

func fetchTorrents(ctx context.Context, url string) ([]*Torrent, error) {
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
			err = se.Unmarshal(&data)
			if err != nil {
				return
			}

			var torrent *Torrent
			torrent, err = data.parse()
			if err != nil {
				return
			}

			torrents = append(torrents, torrent)
		})
		done <- struct{}{}
	})

	go func() {
		visitErr := co.Visit(url)
		if visitErr != nil {
			err = visitErr
		}
		close(done)
	}()

	select {
	case <-done:
		return torrents, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
