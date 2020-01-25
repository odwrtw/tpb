package tpb

import (
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

// UserOptions represent the options for a User lookup
type UserOptions struct {
	OrderBy OrderBy
	Page    int
	Sort    SortOrder
	User    string
}

// Search will search Torrents
func (c *Client) Search(opt SearchOptions) ([]Torrent, error) {
	url := fmt.Sprintf(
		"%s/search/%s/%d/%s/%d/",
		c.Endpoint,
		opt.Key,
		opt.Page,
		mapOrderBy(opt.OrderBy, opt.Sort),
		int(opt.Category),
	)
	return c.parseTorrent(url)
}

// User will return a user's Torrents
func (c *Client) User(opt UserOptions) ([]Torrent, error) {
	url := fmt.Sprintf(
		"%s/user/%s/%d/%s/0/",
		c.Endpoint,
		opt.User,
		opt.Page,
		mapOrderBy(opt.OrderBy, opt.Sort),
	)
	return c.parseTorrent(url)
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

// parseTorrent takes a map of parameters, it will do the request and return
// the parsed Torrents
func (c *Client) parseTorrent(url string) ([]Torrent, error) {
	torrents := []Torrent{}

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
	})

	return torrents, co.Visit(url)
}
