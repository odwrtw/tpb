package tpb

import (
	"fmt"

	"github.com/jpillora/scraper/scraper"
	"github.com/kr/pretty"
	"github.com/mitchellh/mapstructure"
)

// Client represent a Client used to make Search
type Client struct {
	scraper  *scraper.Endpoint
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
	e := &scraper.Endpoint{
		Name:   "TPB",
		Method: "GET",
		List:   "#searchResult > tbody > tr",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2988.133 Safari/537.36",
		},
		Result: map[string]scraper.Extractors{
			"name":     scraper.Extractors{scraper.MustExtractor("a.detLink")},
			"path":     scraper.Extractors{scraper.MustExtractor("a.detLink"), scraper.MustExtractor("@href")},
			"category": scraper.Extractors{scraper.MustExtractor(".vertTh"), scraper.MustExtractor("center"), scraper.MustExtractor("s/\\s//g")},
			"magnet":   scraper.Extractors{scraper.MustExtractor("a[title=Download\\ this\\ torrent\\ using\\ magnet]"), scraper.MustExtractor("@href")},
			"size":     scraper.Extractors{scraper.MustExtractor("/Size (\\d+(\\.\\d+).[KMG]iB)/")},
			"user":     scraper.Extractors{scraper.MustExtractor("a.detDesc")},
			"seeds":    scraper.Extractors{scraper.MustExtractor("td:nth-child(3)")},
			"peers":    scraper.Extractors{scraper.MustExtractor("td:nth-child(4)")},
		},
		Debug: true,
	}

	return &Client{
		scraper:  e,
		Endpoint: endpoint,
	}
}

// OrderBy represent the different values the search can be ordered by
type OrderBy int

// List of the different Orders available
const (
	OrderByName OrderBy = iota
	OrderByDate
	OrderBySize
	OrderBySeeds
	OrderByLeeches
)

// SortOrder represents the sort order
type SortOrder int

// List of the different sort order
const (
	Desc SortOrder = 1 + iota
	Asc
)

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

const defaultSearchPattern = "/search/{{query}}/{{page:0}}/{{orderBy:7}}/{{category:0}}/"
const defaultUserPattern = "/user/{{user}}/{{page:0}}/{{orderBy:7}}/{{category:0}}/"

// Search will search Torrents
func (c *Client) Search(opt SearchOptions) ([]Torrent, error) {
	c.scraper.URL = fmt.Sprintf("%s%s", c.Endpoint, defaultSearchPattern)

	vars := map[string]string{
		"query":    opt.Key,
		"page":     fmt.Sprintf("%d", opt.Page),
		"category": fmt.Sprintf("%d", opt.Category),
		"orderBy":  mapOrderBy(opt.OrderBy, opt.Sort),
	}
	return c.parseTorrent(vars)
}

// User will return a user's Torrents
func (c *Client) User(opt UserOptions) ([]Torrent, error) {
	c.scraper.URL = fmt.Sprintf("%s%s", c.Endpoint, defaultUserPattern)

	vars := map[string]string{
		"user":    opt.User,
		"page":    fmt.Sprintf("%d", opt.Page),
		"orderBy": mapOrderBy(opt.OrderBy, opt.Sort),
	}
	return c.parseTorrent(vars)
}

// mapOrderBy takes the orderBy and sort order parameter and return the
// corresponding option to pass to the website
func mapOrderBy(orderBy OrderBy, order SortOrder) string {
	return fmt.Sprintf("%d", int(orderBy)*2+int(order))
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
func (c *Client) parseTorrent(vars map[string]string) ([]Torrent, error) {
	// Parse the page
	res, err := c.scraper.Execute(vars)
	if err != nil {
		return nil, err
	}

	torrents := []Torrent{}
	pretty.Println(res)

	// Map the res to our structure
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &torrents,
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(res)
	if err != nil {
		return nil, err
	}

	return torrents, nil
}
