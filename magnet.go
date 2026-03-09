package tpb

import (
	"bytes"
	"text/template"
)

var magnetTemplate *template.Template

func init() {
	magnetTemplate = template.Must(template.New("magnet").Parse(magnetTemplateText))
}

const magnetTemplateText = `magnet:?xt=urn:btih:{{.InfoHash}}&dn={{.Name}}{{range .Trackers}}&tr={{.}}{{end}}`

var trackers = []string{
	"udp://tracker.opentrackr.org:1337/announce",
	"udp://open.tracker.cl:1337/announce",
	"udp://tracker.openbittorrent.com:6969/announce",
	"udp://opentracker.i2p.rocks:6969/announce",
	"udp://tracker.torrent.eu.org:451/announce",
	"udp://open.stealth.si:80/announce",
}

// Magnet returns a magnet URI for the torrent. If the API provided a magnet
// link directly, it is returned as-is. Otherwise one is constructed from the
// info hash and a set of well-known trackers.
func (t *Torrent) Magnet() string {
	if t.MagnetLink != "" {
		return t.MagnetLink
	}
	var tpl bytes.Buffer
	tplStruct := struct {
		*Torrent
		Trackers []string
	}{
		Torrent:  t,
		Trackers: trackers,
	}
	if err := magnetTemplate.Execute(&tpl, tplStruct); err != nil {
		return ""
	}
	return tpl.String()
}
