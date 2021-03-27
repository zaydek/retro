package shared

import (
	_ "embed"

	"encoding/json"
)

//go:embed deps.json
var deps string

var Package struct {
	React              string `json:"react"`
	ReactDOM           string `json:"react-dom"`
	Retro              string `json:"@zaydek/retro"`
	RetroStore         string `json:"@zaydek/retro-store"`
	RetroBrowserRouter string `json:"@zaydek/retro-browser-router"`
}

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func init() {
	err := json.Unmarshal([]byte(deps), &Package)
	must(err)
}
