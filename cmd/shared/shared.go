package shared

import (
	_ "embed"

	"encoding/json"
)

//go:embed deps.json
var deps string

type PackageDeps struct {
	ReactVersion    string `json:"react"`
	ReactDOMVersion string `json:"react-dom"`

	TypesReactVersion    string `json:"@types/react"`
	TypesReactDOMVersion string `json:"@types/react-dom"`

	RetroVersion              string `json:"@zaydek/retro"`
	RetroStoreVersion         string `json:"@zaydek/retro-store"`
	RetroBrowserRouterVersion string `json:"@zaydek/retro-browser-router"`
	EsbuildVersion            string `json:"esbuild"`
}

var Deps PackageDeps

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func init() {
	err := json.Unmarshal([]byte(deps), &Deps)
	must(err)
}
