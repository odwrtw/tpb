package tpb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Client represent a Client used to make Search
type Client struct {
	Endpoint string
}

// Torrent represents a Torrent
type Torrent struct {
	Name     string
	Peers    int
	Seeds    int
	User     string
	Magnet   string
	Size     string
	Category string
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
func (c *Client) SearchWithCtx(ctx context.Context, opt SearchOptions) ([]Torrent, error) {
	url := c.Endpoint + opt.String()
	return c.fetchTorrents(ctx, url)
}

// Search will search Torrents
func (c *Client) Search(opt SearchOptions) ([]Torrent, error) {
	ctx := context.Background()
	return c.SearchWithCtx(ctx, opt)
}

// UserWithContext will return a user's Torrents with a context
func (c *Client) UserWithContext(ctx context.Context, opt UserOptions) ([]Torrent, error) {
	url := c.Endpoint + opt.String()
	return c.fetchTorrents(ctx, url)
}

// User will return a user's Torrents
func (c *Client) User(opt UserOptions) ([]Torrent, error) {
	ctx := context.Background()
	return c.UserWithContext(ctx, opt)
}

// FilterByUsers will filter the results and return only those from the given
// users list
func FilterByUsers(torrents []Torrent, users []string) []Torrent {
	filteredTorrents := []Torrent{}
	for _, t := range torrents {
		for _, u := range users {
			if t.User == u {
				filteredTorrents = append(filteredTorrents, t)
			}
		}
	}
	return filteredTorrents
}

func (c *Client) fetchTorrents(ctx context.Context, url string) ([]Torrent, error) {
	torrents := []Torrent{}
	done := make(chan struct{})

	co := colly.NewCollector()
	co.OnHTML("#searchResult", func(e *colly.HTMLElement) {
		e.ForEach("tbody > tr", func(i int, se *colly.HTMLElement) {
			data := struct {
				Name     string `selector:"td > div.detName > a.detLink"`
				Peers    string `selector:"td:nth-child(4)"`
				Seeds    string `selector:"td:nth-child(3)"`
				User     string `selector:"a.detDesc"`
				Magnet   string `selector:"td:nth-child(2) > a:nth-child(2)" attr:"href"`
				Desc     string `selector:"td:nth-child(2) > font"`
				Category string `selector:"td.vertTh > center"`
			}{}
			se.Unmarshal(&data)

			var peers, seeds int
			peers, _ = strconv.Atoi(data.Peers)
			seeds, _ = strconv.Atoi(data.Seeds)

			// Get the size from description
			// e.g. Uploaded 09-10 2011, Size 703.9 MiB, ULed by YIFY
			// the size should be "703.9 MiB"
			var size string
			parts := strings.Fields(data.Desc)
			if len(parts) > 5 {
				size = parts[4] + " " + strings.Trim(parts[5], ",")
			}

			// Parse the category
			category := strings.ReplaceAll(data.Category, "\n", "")
			category = strings.ReplaceAll(category, "\t", "")

			torrents = append(torrents, Torrent{
				Name:     data.Name,
				Peers:    peers,
				Seeds:    seeds,
				User:     data.User,
				Magnet:   data.Magnet,
				Size:     size,
				Category: category,
			})
		})
		done <- struct{}{}
	})

	var err error
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
