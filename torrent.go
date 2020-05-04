package tpb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Torrent represents a Torrent
type Torrent struct {
	ID          int             `json:"id"`
	Category    TorrentCategory `json:"category"`
	Status      UserStatus      `json:"status"`
	Name        string          `json:"name"`
	NumFiles    int             `json:"num_files"`
	Size        uint64          `json:"size"`
	Seeders     int             `json:"seeders"`
	Leechers    int             `json:"leechers"`
	User        string          `json:"username"`
	Added       time.Time       `json:"added"`
	Description string          `json:"descripiton"`
	InfoHash    string          `json:"info_hash"`
	ImdbID      string          `json:"imdb"`
}

// UnmarshalJSON is a custom unmarshal function to handle timestamps and
// boolean as int and convert them to the right type.
func (t *Torrent) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID          flexInt    `json:"id"`
		Category    flexInt    `json:"category"`
		Status      string     `json:"status"`
		Name        string     `json:"name"`
		NumFiles    flexInt    `json:"num_files"`
		InfoHash    string     `json:"info_hash"`
		Description string     `json:"descr"`
		Leechers    flexInt    `json:"leechers"`
		Seeders     flexInt    `json:"seeders"`
		User        string     `json:"username"`
		Size        flexInt    `json:"size"`
		Added       flexInt    `json:"added"`
		ImdbID      flexString `json:"imdb"`
	}

	// Decode json into the aux struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.ID = int(aux.ID)
	t.Category = TorrentCategory(int(aux.Category))
	t.Status = UserStatus(aux.Status)
	t.Name = aux.Name
	t.NumFiles = int(aux.NumFiles)
	t.Size = uint64(aux.Size)
	t.Seeders = int(aux.Seeders)
	t.Leechers = int(aux.Leechers)
	t.User = aux.User
	t.Added = time.Unix(int64(aux.Added), 0)
	t.Description = aux.Description
	t.InfoHash = aux.InfoHash
	t.ImdbID = string(aux.ImdbID)

	return nil
}

// flexInt represents a int that is a string or a int in json
// - if it's a int, it represents the int
// - if it's a string, it represents the string converted in int
type flexInt int

func (fi *flexInt) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("json: Unmarshaling nil data")
	}
	if b[0] != '"' {
		return json.Unmarshal(b, (*int)(fi))
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*fi = flexInt(i)
	return nil
}

// flexString represents a string or 'null' in json
// - if it's a string, it represents the string
// - if it's not, it's ""
type flexString string

func (fs *flexString) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("json: Unmarshaling nil data")
	}
	if b[0] != '"' {
		*fs = ""
		return nil
	}
	return json.Unmarshal(b, (*string)(fs))
}
