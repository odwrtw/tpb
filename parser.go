package tpb

import (
	"errors"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

// Custom parser error
var (
	ErrParserLeechers    = NewParserError("tpb: failed to parse leechers")
	ErrParserSeeders     = NewParserError("tpb: failed to parse seeders")
	ErrParserSize        = NewParserError("tpb: failed to parse size")
	ErrParserCategory    = NewParserError("tpb: failed to parse category")
	ErrParserSubCategory = NewParserError("tpb: failed to parse sub category")
)

// ParserError represents an error due to parsing
type ParserError struct {
	err error
}

// NewParserError returns a new ParserError from an error
func NewParserError(s string) ParserError {
	return ParserError{errors.New(s)}
}

// Unwrap unwraps the error
func (pe ParserError) Unwrap() error {
	return pe.err
}

// Error implements the error interface
func (pe ParserError) Error() string {
	return pe.err.Error()
}

type rawData struct {
	Torrent *Torrent `selector:"-"`

	Name        string `selector:"td > div.detName > a.detLink"`
	Leechers    string `selector:"td:nth-child(4)"`
	Seeders     string `selector:"td:nth-child(3)"`
	User        string `selector:"a.detDesc"`
	Magnet      string `selector:"td:nth-child(2) > a:nth-child(2)" attr:"href"`
	Desc        string `selector:"td:nth-child(2) > font"`
	Category    string `selector:"td.vertTh > center > a:nth-child(1)" attr:"href"`
	SubCategory string `selector:"td.vertTh > center > a:nth-child(3)" attr:"href"`
}

func (rd *rawData) parse() (*Torrent, error) {
	if rd.Torrent == nil {
		rd.Torrent = &Torrent{}
	}

	rd.Torrent.Name = rd.Name
	rd.Torrent.User = rd.User
	rd.Torrent.Magnet = rd.Magnet

	var err error
	rd.Torrent.Leechers, err = strconv.Atoi(rd.Leechers)
	if err != nil {
		return nil, ErrParserLeechers
	}
	rd.Torrent.Seeders, err = strconv.Atoi(rd.Seeders)
	if err != nil {
		return nil, ErrParserSeeders
	}

	// Get the size from description
	// e.g. Uploaded 09-10 2011, Size 703.9 MiB, ULed by YIFY
	// the size should be "703.9 MiB"
	parts := strings.Fields(rd.Desc)
	if len(parts) > 5 {
		rd.Torrent.Size, err = humanize.ParseBytes(parts[4] + " " + strings.Trim(parts[5], ","))
		if err != nil {
			return nil, ErrParserSize
		}
	} else {
		return nil, ErrParserSize
	}

	// Parse the category
	rd.Torrent.Category, err = categoryFromURL(rd.Category)
	if err != nil {
		return nil, ErrParserCategory
	}
	rd.Torrent.SubCategory, err = categoryFromURL(rd.SubCategory)
	if err != nil {
		return nil, ErrParserSubCategory
	}

	return rd.Torrent, nil
}

// URL: /browse/207
func categoryFromURL(url string) (TorrentCategory, error) {
	s := strings.ReplaceAll(url, "/browse/", "")

	v, err := strconv.Atoi(s)
	if err != nil {
		return All, ErrParserCategory
	}

	return TorrentCategory(v), nil
}
