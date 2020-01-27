package tpb

import (
	"errors"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	tt := []struct {
		name            string
		rawData         *rawData
		expectedTorrent *Torrent
		expectedError   error
	}{
		{
			name:          "peers error",
			rawData:       &rawData{Peers: "invalid"},
			expectedError: ErrParserPeers,
		},
		{
			name:          "seeds error",
			rawData:       &rawData{Peers: "0", Seeds: "invalid"},
			expectedError: ErrParserSeeds,
		},
		{
			name: "size error simple",
			rawData: &rawData{
				Peers: "10",
				Seeds: "100",
				Desc:  "invalid",
			},
			expectedError: ErrParserSize,
		},
		{
			name: "size error humanize",
			rawData: &rawData{
				Peers: "10",
				Seeds: "100",
				Desc:  "Uploaded 09-10 2011, Size 703.9 invalid, ULed by YIFY",
			},
			expectedError: ErrParserSize,
		},
		{
			name: "category error",
			rawData: &rawData{
				Peers:    "10",
				Seeds:    "100",
				Desc:     "Uploaded 09-10 2011, Size 703.9 MiB, ULed by YIFY",
				Category: "invalid",
			},
			expectedError: ErrParserCategory,
		},
		{
			name: "sub category error",
			rawData: &rawData{
				Peers:       "10",
				Seeds:       "100",
				Desc:        "Uploaded 09-10 2011, Size 703.9 MiB, ULed by YIFY",
				Category:    "/browse/200",
				SubCategory: "invalid",
			},
			expectedError: ErrParserSubCategory,
		},
		{
			name: "valid parsing",
			rawData: &rawData{
				Name:        "The Matrix",
				User:        "YIFI",
				Magnet:      "magnet://stuff",
				Peers:       "10",
				Seeds:       "100",
				Desc:        "Uploaded 09-10 2011, Size 703.9 MiB, ULed by YIFY",
				Category:    "/browse/200",
				SubCategory: "/browse/207",
			},
			expectedTorrent: &Torrent{
				Name:        "The Matrix",
				User:        "YIFI",
				Magnet:      "magnet://stuff",
				Peers:       10,
				Seeds:       100,
				Size:        738092646,
				Category:    Video,
				SubCategory: VideoHDMovies,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			torrent, err := tc.rawData.parse()
			// Check the error
			if err != tc.expectedError {
				t.Fatalf("expected err %q, got %q", tc.expectedError, err)
			}

			// Check the torrent
			if tc.expectedTorrent == nil {
				return
			}

			if !reflect.DeepEqual(torrent, tc.expectedTorrent) {
				t.Fatalf("expected:\n%+v\ngot:\n%+v", tc.expectedTorrent, torrent)
			}
		})
	}
}

func TestParserError(t *testing.T) {
	err := ErrParserPeers
	if !errors.As(err, &ParserError{}) {
		t.Fatalf("the error should be a parser error")
	}

	errStr := err.Error()
	unwrappedErr := err.Unwrap().Error()
	if errStr != unwrappedErr {
		t.Fatalf("the errors do not match")
	}
}
