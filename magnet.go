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
	"udp://tracker.coppersurfer.tk:6969/announce",
	"udp://9.rarbg.to:2920/announce",
	"udp://tracker.opentrackr.org:1337",
	"udp://tracker.internetwarriors.net:1337/announce",
	"udp://tracker.leechers-paradise.org:6969/announce",
	"udp://tracker.coppersurfer.tk:6969/announce",
	"udp://tracker.pirateparty.gr:6969/announce",
	"udp://tracker.cyberia.is:6969/announce",
}

// Magnet returns a Magnet for a torrent
func (t *Torrent) Magnet() string {
	var tpl bytes.Buffer
	tplStruct := struct {
		*Torrent
		Trackers []string
	}{
		Torrent:  t,
		Trackers: trackers,
	}
	err := magnetTemplate.Execute(&tpl, tplStruct)
	if err != nil {
		return ""
	}
	return tpl.String()
}
