package deps

import (
	_ "embed"
	"fmt"
	"regexp"

	"encoding/json"
)

var (
	//go:embed deps.jsonc
	deps string

	Deps PackageDeps
)

type PackageDeps struct {
	RetroVersion string

	EsbuildVersion string `json:"esbuild"`
	ReactVersion   string `json:"react"`
	SassVersion    string `json:"sass"`
}

var commentsRegex = regexp.MustCompile(`\/\/.*`)

func init() {
	// Remove comments
	deps = commentsRegex.ReplaceAllString(deps, "")
	if err := json.Unmarshal([]byte(deps), &Deps); err != nil {
		panic(fmt.Errorf("json.Unmarshal: %w", err))
	}
}
