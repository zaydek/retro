package retro

import (
	_ "embed"
	"strings"

	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	render "github.com/buildkite/terminal-to-html/v3"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/zaydek/retro/cmd/retro/cli"
	"github.com/zaydek/retro/cmd/retro/pretty"
	"github.com/zaydek/retro/pkg/ipc"
	"github.com/zaydek/retro/pkg/stdio_logger"
	"github.com/zaydek/retro/pkg/terminal"
	"github.com/zaydek/retro/pkg/watch"
)

const (
	WWW_DIR = "www"
	SRC_DIR = "src"
	OUT_DIR = "out"
)

var EPOCH = time.Now()

var (
	cyan    = func(str string) string { return pretty.Accent(str, terminal.Cyan) }
	magenta = func(str string) string { return pretty.Accent(str, terminal.Magenta) }
	red     = func(str string) string { return pretty.Accent(str, terminal.Red) }
)

////////////////////////////////////////////////////////////////////////////////

type BuildResponse struct {
	Errors   []api.Message
	Warnings []api.Message
}

func (r BuildResponse) Dirty() bool {
	return len(r.Errors) > 0 || len(r.Warnings) > 0
}

func (r BuildResponse) String() string {
	errors := api.FormatMessages(r.Errors, api.FormatMessagesOptions{
		Color:         true,
		Kind:          api.ErrorMessage,
		TerminalWidth: 80,
	})
	warnings := api.FormatMessages(r.Warnings, api.FormatMessagesOptions{
		Color:         true,
		Kind:          api.WarningMessage,
		TerminalWidth: 80,
	})
	return strings.Join(append(errors, warnings...), "")
}

func (r BuildResponse) HTML() string {
	str := r.String()
	str = strings.ReplaceAll(str, "╷", "|")
	str = strings.ReplaceAll(str, "│", "|")
	str = strings.ReplaceAll(str, "╵", "|")

	code := string(render.Render([]byte(str)))
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

.term-fg2 { color: #838887; } /* TODO */
.term-fg3 { font-style: italic; }
.term-fg4 { text-decoration: underline; }

/*
 * iTerm ANSI Colors - Normal (Dark Background)
 */
.term-fg30 { color: #000000; }
.term-fg31 { color: #c91b00; }
.term-fg32 { color: #00c200; }
.term-fg33 { color: #c7c400; }
.term-fg34 { color: #0225c7; }
.term-fg35 { color: #c930c7; }
.term-fg36 { color: #00c5c7; }
.term-fg37 { color: #c7c7c7; }

/*
 * iTerm ANSI Colors - Bright (Dark Background)
 */
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
		<pre><code>` + code + `</pre></code>
		<script type="module">const dev = new EventSource("/~dev"); dev.addEventListener("reload", () => { localStorage.setItem("/~dev", "" + Date.now()); window.location.reload() }); dev.addEventListener("error", e => { try { console.error(JSON.parse(e.data)) } catch {} }); window.addEventListener("storage", e => { if (e.key === "/~dev") { window.location.reload() } })</script>
	</body>
</html>
`
}

func (r Runner) Dev() {
	os.Setenv("WWW_DIR", WWW_DIR)
	os.Setenv("SRC_DIR", SRC_DIR)
	os.Setenv("OUT_DIR", OUT_DIR)

	stdin, stdout, stderr, err := ipc.NewCommand("node", "scripts/backend.js")
	if err != nil {
		panic(err)
	}

	dev := make(chan BuildResponse, 1)
	ready := make(chan struct{})

	go func() {
		for result := range watch.Directory(SRC_DIR, 100*time.Millisecond) {
			if result.Error != nil {
				panic(result.Error)
			}
			// // Dedupe stderr
			// if stderrRes.IsDirty() {
			// 	return
			// }
			stdin <- ipc.Request{Kind: "rebuild"}
		}
	}()

	var once sync.Once
	go func() {
		stdin <- ipc.Request{Kind: "build"}
		for {
			select {
			case out := <-stdout:
				once.Do(func() { ready <- struct{}{} })
				var buildRes BuildResponse
				if err := json.Unmarshal(out.Data, &buildRes); err != nil {
					panic(err)
				}
				dev <- buildRes
			case err := <-stderr:
				panic(err)
			}
		}
	}()

	r.Serve(ServerOptions{Dev: dev, Ready: ready})
}

////////////////////////////////////////////////////////////////////////////////

func (r Runner) Build() {
	fmt.Println("TODO")
}

////////////////////////////////////////////////////////////////////////////////

type ServerOptions struct {
	Dev   chan BuildResponse
	Ready chan struct{}
}

func logRequest200(r *http.Request, start time.Time) {
	var durStr string
	if dur := time.Since(start); dur >= time.Millisecond {
		durStr += " "
		durStr += terminal.Dimf("(%s)", pretty.Duration(dur))
	}
	stdio_logger.Stdout(cyan(fmt.Sprintf("'%s %s'%s", r.Method, r.URL.Path, durStr)))
}

func logRequest500(r *http.Request, start time.Time) {
	var durStr string
	if dur := time.Since(start); dur >= time.Millisecond {
		durStr += " "
		durStr += terminal.Dimf("(%s)", pretty.Duration(dur))
	}
	stdio_logger.Stdout(red(fmt.Sprintf("'%s %s'%s", r.Method, r.URL.Path, durStr)))
}

func (r Runner) Serve(opt ServerOptions) {
	var buildRes BuildResponse

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		// 500 Server error
		if buildRes.Dirty() {
			fmt.Fprint(w, buildRes.HTML())
			logRequest500(req, start) // Takes precedence
			// fmt.Fprint(os.Stderr, buildRes)
			return
		}
		// 200 OK - Serve any
		path := getFSPath(req.URL.Path)
		if ext := filepath.Ext(path); ext != "" && ext != ".html" {
			http.ServeFile(w, req, filepath.Join(OUT_DIR, path))
			return
		}
		// 200 OK - Serve index.html
		http.ServeFile(w, req, filepath.Join(OUT_DIR, "index.html"))
		logRequest200(req, start)
	})

	if opt.Dev != nil {
		// Set server-sent event headers
		http.HandleFunc("/~dev", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			flusher, ok := w.(http.Flusher)
			if !ok {
				panic("Internal error")
			}
			for {
				select {
				case buildRes = <-opt.Dev:
					fmt.Fprint(w, "event: reload\ndata\n\n")
					flusher.Flush()
				case <-req.Context().Done():
					return
				}
			}
		})
	}

	if opt.Ready != nil {
		<-opt.Ready
	}

	var durStr string
	if dur := time.Since(EPOCH); dur >= time.Millisecond {
		durStr += " "
		durStr += terminal.Dimf("(%s)", pretty.Duration(time.Since(EPOCH)))
	}

	stdio_logger.Stdout(cyan(fmt.Sprintf("Ready on port '%d'%s", r.getPort(), durStr)))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", r.getPort()), nil); err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////////////////////

var pkg struct {
	React              string `json:"react"`
	ReactDOM           string `json:"react-dom"`
	Retro              string `json:"@zaydek/retro"`
	RetroStore         string `json:"@zaydek/retro-store"`
	RetroBrowserRouter string `json:"@zaydek/retro-browser-router"`
}

//go:embed deps.json
var contents []byte

func Run() {
	if err := json.Unmarshal(contents, &pkg); err != nil {
		panic(err)
	}

	cmd, err := cli.ParseCLIArguments()
	switch err {
	case cli.VersionError:
		fmt.Fprintln(os.Stdout, pkg.Retro)
		os.Exit(0)
	case cli.UsageError:
		fallthrough
	case cli.HelpError:
		fmt.Fprintln(os.Stdout, pretty.Inset(pretty.Spaces(cyan(usage))))
		os.Exit(0)
	}

	switch err.(type) {
	case cli.CommandError:
		fmt.Fprintln(os.Stderr, pretty.Error(magenta(err.Error())))
		os.Exit(1)
	default:
		if err != nil {
			panic(err)
		}
	}

	run := Runner{Command: cmd}
	switch cmd.(type) {
	case cli.DevCommand:
		os.Setenv("NODE_ENV", "development")
		guardErr := guards()
		switch guardErr.(type) {
		case HTMLError:
			fmt.Fprintln(os.Stderr, pretty.Error(magenta(guardErr.Error())))
			os.Exit(1)
		default:
			if guardErr != nil {
				panic(guardErr)
			}
		}
		run.Dev()
	case cli.BuildCommand:
		os.Setenv("NODE_ENV", "production")
		guardErr := guards()
		switch guardErr.(type) {
		case HTMLError:
			fmt.Fprintln(os.Stderr, pretty.Error(magenta(guardErr.Error())))
			os.Exit(1)
		default:
			if guardErr != nil {
				panic(guardErr)
			}
		}
		run.Build()
	case cli.ServeCommand:
		os.Setenv("NODE_ENV", "production")
		run.Serve(ServerOptions{})
	}
}
