package deps

import (
	"testing"

	"github.com/zaydek/retro/go/pkg/expect"
)

func TestVersions(t *testing.T) {
	expect.DeepEqual(t, Deps.EsbuildVersion, "^0.10.2")
	expect.DeepEqual(t, Deps.MDXVersion, "^1.6.22")
	expect.DeepEqual(t, Deps.ReactVersion, "^17.0.2")
	expect.DeepEqual(t, Deps.SassVersion, "^1.32.")
}
