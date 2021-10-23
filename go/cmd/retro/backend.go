package retro

import (
	"path/filepath"
	"strings"

	render "github.com/buildkite/terminal-to-html/v3"
	"github.com/evanw/esbuild/pkg/api"
)

type BundleInfo struct {
	Metafile map[string]interface{}
	Errors   []api.Message
	Warnings []api.Message
}

func (b BundleInfo) IsDirty() bool {
	return len(b.Errors) > 0 || len(b.Warnings) > 0
}

func (b BundleInfo) String() string {
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

func (b BundleInfo) HTML() string {
	renderStr := string(
		render.Render(
			[]byte(
				b.String(),
			),
		),
	)
	return `<!DOCTYPE html>
<html lang="en">
	<head>
		<title>Build Error</title>
		<style>
:root {
	-webkit-font-smoothing: antialiased; /* macOS */
	-moz-osx-font-smoothing: grayscale;  /* Firefox */

	color: #c7c7c7;
	background-color: #000000;
}

code {
	font: 18px / 1.45
		"Monaco",   /* macOS */
		"Consolas", /* Windows */
		monospace;
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
		<pre><code>` + renderStr + `</pre></code>
		` + htmlServerSentEvents + `
	</body>
</html>`
}

////////////////////////////////////////////////////////////////////////////////

type Message struct {
	VendorInfo    BundleInfo
	ClientInfo    BundleInfo
	ClientAppInfo BundleInfo
}

func (m Message) IsDirty() bool {
	return m.VendorInfo.IsDirty() ||
		m.ClientInfo.IsDirty() ||
		m.ClientAppInfo.IsDirty()
}

func (m Message) String() string {
	if m.VendorInfo.IsDirty() {
		return m.VendorInfo.String()
	} else if m.ClientInfo.IsDirty() {
		return m.ClientInfo.String()
	} else if m.ClientAppInfo.IsDirty() {
		return m.ClientAppInfo.String()
	}
	return ""
}

func (m Message) HTML() string {
	if m.VendorInfo.IsDirty() {
		return m.VendorInfo.HTML()
	} else if m.ClientInfo.IsDirty() {
		return m.ClientInfo.HTML()
	} else if m.ClientAppInfo.IsDirty() {
		return m.ClientAppInfo.HTML()
	}
	return ""
}

func (m Message) getChunkedEntrypoints() entryPoints {
	if RETRO_CMD == string(KindDevCommand) {
		if m.VendorInfo.Metafile == nil || m.ClientInfo.Metafile == nil {
			return entryPoints{"client.css", "vendor.js", "client.js"}
		}
	}
	var entries entryPoints
	for key := range m.VendorInfo.Metafile["outputs"].(map[string]interface{}) {
		if strings.HasSuffix(key, ".js") {
			entries.vendorJS, _ = filepath.Rel(RETRO_OUT_DIR, key)
			break
		}
	}
	for key := range m.ClientInfo.Metafile["outputs"].(map[string]interface{}) {
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
