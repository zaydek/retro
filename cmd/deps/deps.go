package deps

import (
	_ "embed"

	"encoding/json"
)

//go:embed deps.json
var deps string

type PackageDeps struct {
	// NOTE: Defer to create_retro_app because 'deps.init' runs before 'main.init'
	RetroVersion string

	ReactVersion         string `json:"react"`
	ReactDOMVersion      string `json:"react-dom"`
	TypesReactVersion    string `json:"@types/react"`
	TypesReactDOMVersion string `json:"@types/react-dom"`
	EsbuildVersion       string `json:"esbuild"`
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
