package retro

import (
	"path/filepath"
	"strings"

	render "github.com/buildkite/terminal-to-html/v3"
	"github.com/evanw/esbuild/pkg/api"
)

type BundleResult struct {
	Metafile map[string]interface{}
	Errors   []api.Message
	Warnings []api.Message
}

func (b BundleResult) IsDirty() bool {
	return len(b.Errors) > 0 || len(b.Warnings) > 0
}

func (b BundleResult) String() string {
	w := api.FormatMessages(b.Warnings, api.FormatMessagesOptions{
		Color:         true,
		Kind:          api.WarningMessage,
		TerminalWidth: 80,
	})
	e := api.FormatMessages(b.Errors, api.FormatMessagesOptions{
		Color:         true,
		Kind:          api.ErrorMessage,
		TerminalWidth: 80,
	})
	return strings.TrimRight(strings.Join(append(e, w...), ""), "\n")
}

func (b BundleResult) HTML() string {
	str := b.String()

	return `<!DOCTYPE html>
<html lang="en">
	<head>
		<title>Build Error</title>
		<style>
html {
	color: #c7c7c7;
	background-color: #000000;
}

code {
	font: 16px / 1.4 "Consolas", "Monaco", monospace;
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
		` + serverSentEventsStub + `
	</body>
</html>`
}

////////////////////////////////////////////////////////////////////////////////

// Describes a `build_done` message or a `rebuild_done` message
type Message struct {
	Kind string
	Data struct {
		Vendor BundleResult
		Client BundleResult
	}
}

func (m Message) GetDirty() BundleResult {
	if m.Data.Vendor.IsDirty() {
		return m.Data.Vendor
	} else if m.Data.Client.IsDirty() {
		return m.Data.Client
	}
	return BundleResult{}
}

func (m Message) getChunkedEntrypoints() entryPoints {
	// Guard zero value messages
	if m.Data.Vendor.Metafile == nil || m.Data.Client.Metafile == nil {
		return entryPoints{}
	}
	var entries entryPoints
	for key := range m.Data.Vendor.Metafile["outputs"].(map[string]interface{}) {
		if strings.HasSuffix(key, ".js") {
			entries.vendorJS, _ = filepath.Rel(RETRO_OUT_DIR, key)
			break
		}
	}
	for key := range m.Data.Client.Metafile["outputs"].(map[string]interface{}) {
		if strings.HasSuffix(key, ".css") {
			entries.clientCSS, _ = filepath.Rel(RETRO_OUT_DIR, key)
			if entries.clientJS != "" { // Check other
				break
			}
		} else if strings.HasSuffix(key, ".js") {
			entries.clientJS, _ = filepath.Rel(RETRO_OUT_DIR, key)
			if entries.clientCSS != "" { // Check other
				break
			}
		}
	}
	return entries
}
