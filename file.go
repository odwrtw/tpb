package tpb

import (
	"encoding/json"
	"path"
)

// File represents a file within a torrent
type File struct {
	// Name is the full path of the file, joined from the path components
	// returned by the API
	Name string
	// Size is the file size in bytes
	Size int64
}

// UnmarshalJSON handles the API's format where name is an array of path
// components and size is a string
func (f *File) UnmarshalJSON(data []byte) error {
	var aux struct {
		Name []string `json:"name"`
		Size flexInt  `json:"size"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	f.Name = path.Join(aux.Name...)
	f.Size = int64(aux.Size)
	return nil
}
