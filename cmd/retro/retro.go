package retro

import (
	_ "embed"
	"io/ioutil"
	"regexp"
	"sort"
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

type BackendResponse struct {
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

func (r Runner) Dev() {
	stdin, stdout, stderr, err := ipc.NewCommand("node", "scripts/backend.esbuild.js")
	if err != nil {
		panic(err)
	}

	dev := make(chan BackendResponse, 1)
	ready := make(chan struct{})

	go func() {
		for result := range watch.Directory(SRC_DIR, 100*time.Millisecond) {
			if result.Error != nil {
				panic(result.Error)
			}
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
				var res BackendResponse
				if err := json.Unmarshal(out.Data, &res); err != nil {
					panic(err)
				}
				dev <- res
			case err := <-stderr:
				panic(err)
			}
		}
	}()

	r.Serve(ServerOptions{Dev: dev, Ready: ready})
}

////////////////////////////////////////////////////////////////////////////////

type lsInfo struct {
	path string
	size int64
}

type lsInfos []lsInfo

func (a lsInfos) Len() int           { return len(a) }
func (a lsInfos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a lsInfos) Less(i, j int) bool { return a[i].path < a[j].path }

func ls(dir string) (lsInfos, error) {
	var ls lsInfos
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ls = append(ls, lsInfo{
			path: path,
			size: info.Size(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ls, nil
}

// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format
func byteCount(b int64) string {
	const u = 1024

	if b < u {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(u), 0
	for n := b / u; n >= u; n /= u {
		div *= u
		exp++
	}
	return fmt.Sprintf("%.0f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

var greedyExtRe = regexp.MustCompile(`(\.).*$`)

func (r Runner) Build() {
	stdin, stdout, stderr, err := ipc.NewCommand("node", "scripts/backend.esbuild.js")
	if err != nil {
		panic(err)
	}

	stdin <- ipc.Request{Kind: "build"}
	select {
	case out := <-stdout:
		var res BackendResponse
		if err := json.Unmarshal(out.Data, &res); err != nil {
			panic(err)
		}
		if !res.Dirty() {
			break
		}
		fmt.Fprint(os.Stderr, res) // Use fmt.Fprint not fmt.Fprintln
		os.Exit(1)
	case err := <-stderr:
		panic(err)
	}

	infos, err := ls(OUT_DIR)
	if err != nil {
		panic(err)
	}

	sort.Sort(infos)

	var sum, sumMap int64
	for x, v := range infos {
		var hue = terminal.Normal
		if strings.HasSuffix(v.path, ".html") {
			hue = terminal.Normal
		} else if strings.HasSuffix(v.path, ".js") || strings.HasSuffix(v.path, ".js.map") {
			hue = terminal.Yellow
		} else if strings.HasSuffix(v.path, ".css") || strings.HasSuffix(v.path, ".css.map") {
			hue = terminal.Cyan
		} else {
			hue = terminal.Dim
		}

		if x == 0 {
			fmt.Println()
		}
		fmt.Printf(" %v%s%v\n",
			hue(v.path),
			strings.Repeat(" ", 25-len(v.path)),
			terminal.Dimf("(%s)", byteCount(v.size)),
		)

		if !strings.HasSuffix(v.path, ".map") {
			sum += v.size
		}
		sumMap += v.size
	}

	fmt.Println()
	fmt.Println(strings.Repeat(" ", 25), terminal.Dimf("(%s)", byteCount(sum)))
	fmt.Println(strings.Repeat(" ", 25), terminal.Dimf("(%s w/ sourcemaps)", byteCount(sumMap)))

	durStr := terminal.Dimf("(%s)", pretty.Duration(time.Since(EPOCH)))

	fmt.Println()
	fmt.Println(fmt.Sprintf("%s", durStr))

}

////////////////////////////////////////////////////////////////////////////////

type ServerOptions struct {
	Dev   chan BackendResponse
	Ready chan struct{}
}

func serve200Str(r *http.Request, start time.Time) string {
	var durStr string
	if dur := time.Since(start); dur >= time.Millisecond {
		durStr += " "
		durStr += terminal.Dimf("(%s)", pretty.Duration(dur))
	}
	return cyan(fmt.Sprintf("'%s %s'%s", r.Method, r.URL.Path, durStr))
}

func serve500Str(r *http.Request, start time.Time) string {
	var durStr string
	if dur := time.Since(start); dur >= time.Millisecond {
		durStr += " "
		durStr += terminal.Dimf("(%s)", pretty.Duration(dur))
	}
	return red(fmt.Sprintf("'%s %s'%s", r.Method, r.URL.Path, durStr))
}

func (r Runner) Serve(opt ServerOptions) {
	var res BackendResponse

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// 500 Server error
		start := time.Now()
		if res.Dirty() {
			fmt.Fprint(w, res.HTML())
			stdio_logger.Stdout(serve500Str(req, start))
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
		stdio_logger.Stdout(serve200Str(req, start))
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
				case res = <-opt.Dev:
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
var deps string

// Server-sent events and localStorage events
const devStub = `<script type="module">const dev = new EventSource("/~dev"); dev.addEventListener("reload", () => { localStorage.setItem("/~dev", "" + Date.now()); window.location.reload() }); dev.addEventListener("error", e => { try { console.error(JSON.parse(e.data)) } catch {} }); window.addEventListener("storage", e => { if (e.key === "/~dev") { window.location.reload() } })</script>`

func copyEntryPoint() error {
	bstr, err := ioutil.ReadFile(filepath.Join(WWW_DIR, "index.html"))
	if err != nil {
		return err
	}
	contents := string(bstr)
	if os.Getenv("ENV") == "development" {
		contents = strings.Replace(contents, ">\n</html>", fmt.Sprintf(">\n\t%s\n</html>", devStub), 1)
	}
	if err := ioutil.WriteFile(filepath.Join(OUT_DIR, "index.html"), []byte(contents), MODE_FILE); err != nil {
		return err
	}
	return nil
}

func Run() {
	if err := json.Unmarshal([]byte(deps), &pkg); err != nil {
		panic(err)
	}

	cmd, err := cli.ParseCLIArguments()
	switch err {
	case cli.VersionError:
		fmt.Println(pkg.Retro)
		return
	case cli.UsageError:
		fallthrough
	case cli.HelpError:
		fmt.Println(pretty.Inset(pretty.Spaces(cyan(usage))))
		return
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

		os.Setenv("CMD", "dev")
		os.Setenv("ENV", "development")
		os.Setenv("WWW_DIR", WWW_DIR)
		os.Setenv("SRC_DIR", SRC_DIR)
		os.Setenv("OUT_DIR", OUT_DIR)

		rmdirs := []string{OUT_DIR}
		for _, rmdir := range rmdirs {
			if err := os.RemoveAll(rmdir); err != nil {
				panic(err)
			}
		}

		mkdirs := []string{WWW_DIR, SRC_DIR, OUT_DIR}
		for _, mkdir := range mkdirs {
			if err := os.MkdirAll(mkdir, MODE_DIR); err != nil {
				panic(err)
			}
		}

		guardErr := entryPointGuards()
		switch guardErr.(type) {
		case HTMLError:
			fmt.Fprintln(os.Stderr, pretty.Error(magenta(guardErr.Error())))
			os.Exit(1)
		default:
			if guardErr != nil {
				panic(guardErr)
			}
		}

		if err := copyAll(WWW_DIR, filepath.Join(OUT_DIR, WWW_DIR), []string{"index.html"}); err != nil {
			panic(err)
		}
		if err := copyEntryPoint(); err != nil {
			panic(err)
		}

		run.Dev()
	case cli.BuildCommand:

		os.Setenv("CMD", "build")
		os.Setenv("ENV", "production")
		os.Setenv("WWW_DIR", WWW_DIR)
		os.Setenv("SRC_DIR", SRC_DIR)
		os.Setenv("OUT_DIR", OUT_DIR)

		rmdirs := []string{OUT_DIR}
		for _, rmdir := range rmdirs {
			if err := os.RemoveAll(rmdir); err != nil {
				panic(err)
			}
		}

		mkdirs := []string{WWW_DIR, SRC_DIR, OUT_DIR}
		for _, mkdir := range mkdirs {
			if err := os.MkdirAll(mkdir, MODE_DIR); err != nil {
				panic(err)
			}
		}

		guardErr := entryPointGuards()
		switch guardErr.(type) {
		case HTMLError:
			fmt.Fprintln(os.Stderr, pretty.Error(magenta(guardErr.Error())))
			os.Exit(1)
		default:
			if guardErr != nil {
				panic(guardErr)
			}
		}

		if err := copyAll(WWW_DIR, filepath.Join(OUT_DIR, WWW_DIR), []string{filepath.Join(WWW_DIR, "index.html")}); err != nil {
			panic(err)
		}
		if err := copyEntryPoint(); err != nil {
			panic(err)
		}

		run.Build()
	case cli.ServeCommand:

		os.Setenv("CMD", "serve")
		os.Setenv("ENV", "production")
		os.Setenv("WWW_DIR", WWW_DIR)
		os.Setenv("SRC_DIR", SRC_DIR)
		os.Setenv("OUT_DIR", OUT_DIR)

		run.Serve(ServerOptions{})
	}
}
