package tpb

// TorrentCategory represents the torrent categories
type TorrentCategory int

// List of all the available categories
const (
	// All is all the results
	All TorrentCategory = 0

	// Audio matches only audio results
	Audio           = 100
	AudioMusic      = 101
	AudioAudiobooks = 102
	AudioSoundclips = 103
	AudioFLAC       = 104
	AudioOther      = 199

	// Video matches only video results
	Video            = 200
	VideoMovies      = 201
	VideoMoviesDVDR  = 202
	VideoMusicvideos = 203
	VideoMovieclips  = 204
	VideoTVshows     = 205
	VideoHandheld    = 206
	VideoHDMovies    = 207
	VideoHDTVshows   = 208
	Video3D          = 209
	VideoOther       = 299

	// Applications matches only app results
	Applications         = 300
	ApplicationsWindows  = 301
	ApplicationsMac      = 302
	ApplicationsUNIX     = 303
	ApplicationsHandheld = 304
	ApplicationsIOS      = 305
	ApplicationsAndroid  = 306
	ApplicationsOtherOS  = 399

	// Games matches only game results
	Games         = 400
	GamesPC       = 401
	GamesMac      = 402
	GamesPSx      = 403
	GamesXBOX360  = 404
	GamesWii      = 405
	GamesHandheld = 406
	GamesIOS      = 407
	GamesAndroid  = 408
	GamesOther    = 499

	// Other matches all the other categories
	Other          = 600
	OtherEbooks    = 601
	OtherComics    = 602
	OtherPictures  = 603
	OtherCovers    = 604
	OtherPhysibles = 605
	OtherOther     = 699
)
