package retro

import (
	_ "embed"
	"io/ioutil"
	"sort"
	"strings"

	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/zaydek/retro/cmd/retro/cli"
	"github.com/zaydek/retro/cmd/retro/pretty"
	"github.com/zaydek/retro/cmd/shared"
	"github.com/zaydek/retro/pkg/ipc"
	"github.com/zaydek/retro/pkg/stdio_logger"
	"github.com/zaydek/retro/pkg/terminal"
	"github.com/zaydek/retro/pkg/watch"
)

var EPOCH = time.Now()

var (
	cyan    = func(str string) string { return pretty.Accent(str, terminal.Cyan) }
	magenta = func(str string) string { return pretty.Accent(str, terminal.Magenta) }
	red     = func(str string) string { return pretty.Accent(str, terminal.Red) }
)

////////////////////////////////////////////////////////////////////////////////

type DevOptions struct {
	Preflight bool
}

func (r Runner) Dev(opt DevOptions) {
	var copyHTMLEntryPoint func(string, string, string) error
	if opt.Preflight {
		var err error
		copyHTMLEntryPoint, err = r.preflight()
		switch err.(type) {
		case HTMLError:
			fmt.Fprintln(os.Stderr, pretty.Error(magenta(err.Error())))
			os.Exit(1)
		default:
			if err != nil {
				panic(err)
			}
		}
	}

	stdin, stdout, stderr, err := ipc.NewCommand("node", "scripts/backend.esbuild.js")
	if err != nil {
		panic(err)
	}

	dev := make(chan BackendResponse, 1)
	ready := make(chan struct{})

	go func() {
		for result := range watch.Directory(SRC_DIR, 100*time.Millisecond) {
			if result.Err != nil {
				panic(result.Err)
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
				var res BackendResponse
				if err := json.Unmarshal(out.Data, &res); err != nil {
					panic(err)
				}
				once.Do(func() {
					if err := copyHTMLEntryPoint("react.js", "index.js", "index.css"); err != nil {
						panic(err)
					}
					ready <- struct{}{}
				})
				dev <- res
			case err := <-stderr:
				panic(err)
			}
		}
	}()

	r.Serve(ServerOptions{Dev: dev, Ready: ready})
}

////////////////////////////////////////////////////////////////////////////////

type BuildOptions struct {
	Preflight bool
}

func (r Runner) Build(opt BuildOptions) {
	var copyHTMLEntryPoint func(string, string, string) error
	if opt.Preflight {
		var err error
		copyHTMLEntryPoint, err = r.preflight()
		switch err.(type) {
		case HTMLError:
			fmt.Fprintln(os.Stderr, pretty.Error(magenta(err.Error())))
			os.Exit(1)
		default:
			if err != nil {
				panic(err)
			}
		}
	}

	stdin, stdout, stderr, err := ipc.NewCommand("node", "scripts/backend.esbuild.js")
	if err != nil {
		panic(err)
	}

	stdin <- ipc.Request{Kind: "build"}

	var once sync.Once
	select {
	case out := <-stdout:
		var res BackendResponse
		if err := json.Unmarshal(out.Data, &res); err != nil {
			panic(err)
		}
		once.Do(func() {
			react_js, index_js, index_css := res.getChunkedNames()
			if err := copyHTMLEntryPoint(react_js, index_js, index_css); err != nil {
				panic(err)
			}
		})
		if res.Dirty() {
			fmt.Fprint(os.Stderr, res)
			os.Exit(1)
		}
	case err := <-stderr:
		panic(err)
	}

	infos, err := ls(OUT_DIR)
	if err != nil {
		panic(err)
	}

	sort.Sort(infos)

	var sum, sumMap int64
	for _, v := range infos {
		var color = terminal.Normal
		if strings.HasSuffix(v.path, ".html") {
			color = terminal.Normal
		} else if strings.HasSuffix(v.path, ".js") || strings.HasSuffix(v.path, ".js.map") {
			color = terminal.Yellow
		} else if strings.HasSuffix(v.path, ".css") || strings.HasSuffix(v.path, ".css.map") {
			color = terminal.Cyan
		} else {
			color = terminal.Dim
		}

		fmt.Printf("%v%s%v\n",
			color(v.path),
			strings.Repeat(" ", 32-len(v.path)),
			terminal.Dimf("(%s)", byteCount(v.size)),
		)

		if !strings.HasSuffix(v.path, ".map") {
			sum += v.size
		}
		sumMap += v.size
	}

	// TODO: Wrap w/ 'if r.Sourcemap { ... }'
	fmt.Println(strings.Repeat(" ", 32) + terminal.Dimf("(%s sum)", byteCount(sum)))
	fmt.Println(strings.Repeat(" ", 32) + terminal.Dimf("(%s sum w/ sourcemaps)", byteCount(sumMap)))

	durStr := terminal.Dimf("(%s)", pretty.Duration(time.Since(EPOCH)))

	fmt.Println()
	fmt.Println(fmt.Sprintf("%s", durStr))

}

////////////////////////////////////////////////////////////////////////////////

type ServerOptions struct {
	Preflight bool

	Dev   chan BackendResponse
	Ready chan struct{}
}

func formatServe200(r *http.Request, start time.Time) string {
	var durStr string
	if dur := time.Since(start); dur >= time.Millisecond {
		durStr += " "
		durStr += terminal.Dimf("(%s)", pretty.Duration(dur))
	}
	return cyan(fmt.Sprintf("'%s %s'%s", r.Method, r.URL.Path, durStr))
}

func formatServe500(r *http.Request, start time.Time) string {
	var durStr string
	if dur := time.Since(start); dur >= time.Millisecond {
		durStr += " "
		durStr += terminal.Dimf("(%s)", pretty.Duration(dur))
	}
	return red(fmt.Sprintf("'%s %s'%s", r.Method, r.URL.Path, durStr))
}

func (r Runner) Serve(opt ServerOptions) {
	if opt.Preflight {
		_, err := r.preflight()
		switch err.(type) {
		case HTMLError:
			fmt.Fprintln(os.Stderr, pretty.Error(magenta(err.Error())))
			os.Exit(1)
		default:
			if err != nil {
				panic(err)
			}
		}
	}

	if opt.Ready != nil {
		<-opt.Ready
	}

	var res BackendResponse

	// Add the dev stub
	var contents string
	if os.Getenv("ENV") == "development" {
		bstr, err := ioutil.ReadFile(filepath.Join(OUT_DIR, "index.html"))
		if err != nil {
			panic(err)
		}
		contents = string(bstr)
		contents = strings.Replace(
			contents,
			"</body>",
			fmt.Sprintf("\t%s\n\t</body>", fmt.Sprintf(`<script type="module">%s</script>`, devStub)),
			1,
		)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/~dev" {
			return
		}

		// 500 Server error
		start := time.Now()
		if res.Dirty() {
			fmt.Fprint(w, res.HTML())
			stdio_logger.Stdout(formatServe500(req, start))
			return
		}
		// 200 OK - Serve any
		path := getFSPath(req.URL.Path)
		if ext := filepath.Ext(path); ext != "" && ext != ".html" {
			http.ServeFile(w, req, filepath.Join(OUT_DIR, path))
			return
		}
		// 200 OK - Serve index.html
		if os.Getenv("ENV") == "development" {
			fmt.Fprint(w, contents)
		} else {
			http.ServeFile(w, req, filepath.Join(OUT_DIR, "index.html"))
		}
		stdio_logger.Stdout(formatServe200(req, start))
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

func Run() {
	cmd, err := cli.ParseCLIArguments()
	switch err {
	case cli.VersionError:
		fmt.Println(shared.Deps.RetroVersion)
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
		run.Dev(DevOptions{Preflight: true})
	case cli.BuildCommand:
		run.Build(BuildOptions{Preflight: true})
	case cli.ServeCommand:
		run.Serve(ServerOptions{Preflight: true})
	}
}
