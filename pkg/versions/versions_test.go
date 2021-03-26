package versions

import (
	_ "embed"
	"testing"

	"github.com/zaydek/retro/pkg/expect"
)

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func TestExact(t *testing.T) {
	contents := `
+---------------------+
| react     | v17.0.1 |
| react-dom | v17.0.1 |
+---------------------+
`

	vs, err := Parse(contents)
	must(err)
	expect.DeepEqual(t, vs, Versions{
		"react":     "v17.0.1",
		"react-dom": "v17.0.1",
	})
}

func TestLatest(t *testing.T) {
	contents := `
+--------------------+
| react     | latest |
| react-dom | latest |
+--------------------+
`

	vs, err := Parse(contents)
	must(err)
	expect.DeepEqual(t, vs, Versions{
		"react":     "latest",
		"react-dom": "latest",
	})
}
