package tpb

import (
	"context"
	"math/rand"
	"time"

	"github.com/gocolly/colly/v2"
)

type ua []string

func (u ua) random() string {
	rand.Seed(time.Now().UnixNano())
	return u[rand.Intn(len(u))]
}

var userAgents = ua{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
	"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1",
}

func fetchTorrents(ctx context.Context, url string) ([]*Torrent, error) {
	var err error
	torrents := []*Torrent{}
	done := make(chan struct{})

	co := colly.NewCollector(
		colly.UserAgent(userAgents.random()),
	)
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
