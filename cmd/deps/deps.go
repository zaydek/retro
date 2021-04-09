package deps

import (
	_ "embed"
	"regexp"

	"encoding/json"
)

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

var (
	//go:embed deps.jsonc
	deps string

	Deps PackageDeps
)

type PackageDeps struct {
	RetroVersion string

	EsbuildVersion string `json:"esbuild"`
	MDXVersion     string `json:"mdx"`
	ReactVersion   string `json:"react"`
	SassVersion    string `json:"sass"`
}

var re = regexp.MustCompile(`\/\/.*`)

func init() {
	// Remove comments
	deps = re.ReplaceAllString(deps, "")
	err := json.Unmarshal([]byte(deps), &Deps)
	must(err)
}
