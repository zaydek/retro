package retro

import (
	"path/filepath"
	"strings"

	render "github.com/buildkite/terminal-to-html/v3"
	"github.com/evanw/esbuild/pkg/api"
)

type BackendResponse struct {
	Metafile struct {
		Vendor map[string]interface{}
		Bundle map[string]interface{}
	}
	Errors   []api.Message
	Warnings []api.Message
}

func (r BackendResponse) Dirty() bool {
	return len(r.Errors) > 0 || len(r.Warnings) > 0
}

func (r BackendResponse) String() string {
	e := api.FormatMessages(r.Errors, api.FormatMessagesOptions{
		Color:         true,
		Kind:          api.ErrorMessage,
		TerminalWidth: 80,
	})
	w := api.FormatMessages(r.Warnings, api.FormatMessagesOptions{
		Color:         true,
		Kind:          api.WarningMessage,
		TerminalWidth: 80,
	})
	return strings.Join(append(e, w...), "")
}

func (r BackendResponse) HTML() string {
	str := r.String()
	str = strings.ReplaceAll(str, "╷", "|")
	str = strings.ReplaceAll(str, "│", "|")
	str = strings.ReplaceAll(str, "╵", "|")

	return `<!DOCTYPE html>
<html>
	<head>
		<title>Error</title>
		<style>

html {
	color: #c7c7c7;
	background-color: #000000;
}

code {
	font: 16px / 1.4 "Monaco", monospace;
}

a { color: unset; text-decoration: unset; }
a:hover { text-decoration: underline; }

.term-fg1 {
	font-weight: bold;
	color: #feffff;
}

.term-fg30 { color: #000000; }
.term-fg31 { color: #c91b00; }
.term-fg32 { color: #00c200; }
.term-fg33 { color: #c7c400; }
.term-fg34 { color: #0225c7; }
.term-fg35 { color: #c930c7; }
.term-fg36 { color: #00c5c7; }
.term-fg37 { color: #c7c7c7; }

.term-fg1.term-fg30 { color: #676767; }
.term-fg1.term-fg31 { color: #ff6d67; }
.term-fg1.term-fg32 { color: #5ff967; }
.term-fg1.term-fg33 { color: #fefb67; }
.term-fg1.term-fg34 { color: #6871ff; }
.term-fg1.term-fg35 { color: #ff76ff; }
.term-fg1.term-fg36 { color: #5ffdff; }
.term-fg1.term-fg37 { color: #feffff; }

		</style>
	</head>
	<body>
		<pre><code>` + string(render.Render([]byte(str))) + `</pre></code>
	</body>
	` + devStub + `
</html>
`
}

func (r BackendResponse) getChunkedNames() (react_js, index_js, index_css string) {
	for k := range r.Metafile.Vendor["outputs"].(map[string]interface{}) {
		if strings.HasSuffix(k, ".js") {
			react_js, _ = filepath.Rel(OUT_DIR, k)
			break
		}
	}
	for k := range r.Metafile.Bundle["outputs"].(map[string]interface{}) {
		if strings.HasSuffix(k, ".js") {
			index_js, _ = filepath.Rel(OUT_DIR, k)
			if index_css != "" { // Check other
				break
			}
		} else if strings.HasSuffix(k, ".css") {
			index_css, _ = filepath.Rel(OUT_DIR, k)
			if index_js != "" { // Check other
				break
			}
		}
	}
	return
}
