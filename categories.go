package tpb

import (
	"net/url"
	"strconv"
)

// TorrentCategory represents the torrent categories
type TorrentCategory int

// List of all the available categories
const (
	// All is all the results
	All TorrentCategory = 0

	// Audio matches only audio results
	Audio           TorrentCategory = 100
	AudioMusic      TorrentCategory = 101
	AudioAudiobooks TorrentCategory = 102
	AudioSoundclips TorrentCategory = 103
	AudioFLAC       TorrentCategory = 104
	AudioOther      TorrentCategory = 199

	// Video matches only video results
	Video            TorrentCategory = 200
	VideoMovies      TorrentCategory = 201
	VideoMoviesDVDR  TorrentCategory = 202
	VideoMusicvideos TorrentCategory = 203
	VideoMovieclips  TorrentCategory = 204
	VideoTVshows     TorrentCategory = 205
	VideoHandheld    TorrentCategory = 206
	VideoHDMovies    TorrentCategory = 207
	VideoHDTVshows   TorrentCategory = 208
	Video3D          TorrentCategory = 209
	VideoOther       TorrentCategory = 299

	// Applications matches only app results
	Applications         TorrentCategory = 300
	ApplicationsWindows  TorrentCategory = 301
	ApplicationsMac      TorrentCategory = 302
	ApplicationsUNIX     TorrentCategory = 303
	ApplicationsHandheld TorrentCategory = 304
	ApplicationsIOS      TorrentCategory = 305
	ApplicationsAndroid  TorrentCategory = 306
	ApplicationsOtherOS  TorrentCategory = 399

	// Games matches only game results
	Games         TorrentCategory = 400
	GamesPC       TorrentCategory = 401
	GamesMac      TorrentCategory = 402
	GamesPSx      TorrentCategory = 403
	GamesXBOX360  TorrentCategory = 404
	GamesWii      TorrentCategory = 405
	GamesHandheld TorrentCategory = 406
	GamesIOS      TorrentCategory = 407
	GamesAndroid  TorrentCategory = 408
	GamesOther    TorrentCategory = 499

	// Other matches all the other categories
	Other          TorrentCategory = 600
	OtherEbooks    TorrentCategory = 601
	OtherComics    TorrentCategory = 602
	OtherPictures  TorrentCategory = 603
	OtherCovers    TorrentCategory = 604
	OtherPhysibles TorrentCategory = 605
	OtherOther     TorrentCategory = 699
)

// URLValue returns the url.Values for a TorrentCategory
func (cat *TorrentCategory) URLValue() *url.Values {
	v := url.Values{}
	v.Add("cat", strconv.Itoa(int(*cat)))
	return &v
}
