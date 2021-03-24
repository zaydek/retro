package retro

import (
	"fmt"
	"html"
	"os"
	"strconv"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

type BuildResponse struct {
	Errors   []api.Message
	Warnings []api.Message
}

func (r BuildResponse) Dirty() bool {
	return len(r.Errors) > 0 || len(r.Warnings) > 0
}

type MessageKind int

const (
	Error MessageKind = iota
	Warning
)

func getPos(x interface{}, offset int) string {
	switch y := x.(type) {
	case api.Message:
		return fmt.Sprintf("%s:%d:%d", y.Location.File, y.Location.Line, y.Location.Column+offset)
	case api.Note:
		return fmt.Sprintf("%s:%d:%d", y.Location.File, y.Location.Line, y.Location.Column+offset)
	}
	panic("Internal error")
}

func getVSCodePos(x interface{}) string {
	cwd, _ := os.Getwd()
	switch y := x.(type) {
	case api.Message:
		return fmt.Sprintf("vscode://file%s/%s", cwd, getPos(y, 1))
	case api.Note:
		return fmt.Sprintf("vscode://file%s/%s", cwd, getPos(y, 1))
	}
	panic("Internal error")
}

func pipes(str string) string {
	out := str
	out = strings.Replace(out, "╵", "|", -1)
	out = strings.Replace(out, "│", "|", -1)
	out = strings.Replace(out, "╷", "|", -1)
	return out
}

func formatMessage(m api.Message, kind MessageKind) string {
	var class, typ string
	switch kind {
	case Error:
		class = "red"
		typ = "error"
	case Warning:
		class = "yellow"
		typ = "warning"
	}

	// lineText := strings.ReplaceAll(m.Location.LineText, "\t", "  ")

	tabLen := func(str string) int {
		var len int
		for _, ch := range str {
			if ch == '\t' {
				len += 2
				continue
			}
			len++
		}
		return len
	}

	focus := "^"
	if m.Location.Length > 0 { // FIXME
		focus = strings.Repeat("~", m.Location.Length)
	}

	text := strings.Split(m.Location.LineText, "\n")[0]

	var str string
	str += fmt.Sprintf(`<strong class="bold"> &gt; <a href="%s">%s</a>: <span class="%s">%s:</span> %s</strong>
    %d │ %s<span class="focus">%s</span>%s
    %s │ %s<span class="focus">%s</span>
`,
		getVSCodePos(m),
		getPos(m, 0),
		class,
		typ,
		pipes(m.Text),
		m.Location.Line,
		html.EscapeString(text[:m.Location.Column]),
		html.EscapeString(text[m.Location.Column:m.Location.Column+m.Location.Length]),
		html.EscapeString(text[m.Location.Column+m.Location.Length:]),
		strings.Repeat(" ", len(strconv.Itoa(m.Location.Line))),
		strings.Repeat(" ", tabLen(text[:m.Location.Column])),
		focus,
	)
	if len(m.Notes) > 0 {
		for _, n := range m.Notes {

			focus := "^"
			if n.Location.Length > 0 { // FIXME
				focus = strings.Repeat("~", n.Location.Length)
			}

			text := strings.Split(n.Location.LineText, "\n")[0]

			str += fmt.Sprintf(`   <a href="%s">%s</a>: <span class="bold">note:</span> %s
    %d │ %s<span class="focus">%s</span>%s
    %s │ %s<span class="focus">%s</span>
`,
				getVSCodePos(n),
				getPos(n, 0),
				pipes(n.Text),
				n.Location.Line,
				html.EscapeString(text[:n.Location.Column]),
				html.EscapeString(text[n.Location.Column:n.Location.Column+n.Location.Length]),
				html.EscapeString(text[n.Location.Column+n.Location.Length:]),
				strings.Repeat(" ", len(strconv.Itoa(n.Location.Line))),
				strings.Repeat(" ", tabLen(text[:n.Location.Column])),
				focus,
			)
		}
	}
	str = strings.ReplaceAll(str, "\t", "  ")
	return str
}

func (r BuildResponse) HTML() string {
	var body string

	for _, msg := range r.Errors {
		if body != "" {
			body += "<br>"
		}
		body += `<pre><code>` + formatMessage(msg, Error) + `</code></pre>`
	}

	for _, msg := range r.Warnings {
		if body != "" {
			body += "<br>"
		}
		body += `<pre><code>` + formatMessage(msg, Warning) + `</code></pre>`
	}

	return `<!DOCTYPE html>
<html>
	<head>
		<title>` + fmt.Sprintf("Error: %s", r.Errors[0].Text) + `</title>
		<style>

*,
*::before,
*::after {
	/* Zero-out margin and padding */
	margin: 0;
	padding: 0;

	/* Reset border-box */
  box-sizing: border-box;
}

:root {
	--font-size: 14px;

	--color: #c7c7c7;
	--bg: #000000;

	--bold: #feffff;
	--red: #ff6d67;
	--yellow: #fefb67;
	--focus: #00c200;
}

body {
  color: #c7c7c7;
  background-color: #000000;
}

.terminal {
	padding: var(--font-size);
}

code {
	font: var(--font-size) / calc(1.5 * var(--font-size)) "Monaco", monospace;
}

a { color: unset; text-decoration: unset; }
a:hover { text-decoration: underline; }

.bold { color: var(--bold); }
.red { color: var(--red); }
.yellow { color: var(--yellow); }
.focus { color: var(--focus); }

		</style>
	</head>
	<body>
		<div class="terminal">` + body + `</div>
		<script type="module">const dev = new EventSource("/~dev"); dev.addEventListener("reload", () => { localStorage.setItem("/~dev", "" + Date.now()); window.location.reload() }); dev.addEventListener("error", e => { try { console.error(JSON.parse(e.data)) } catch {} }); window.addEventListener("storage", e => { if (e.key === "/~dev") { window.location.reload() } })</script>
	</body>
</html>
`
}
