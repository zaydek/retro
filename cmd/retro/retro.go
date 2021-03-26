package retro

import (
	_ "embed"

	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/zaydek/retro/cmd/retro/cli"
	"github.com/zaydek/retro/cmd/retro/pretty"
	"github.com/zaydek/retro/pkg/ipc"
	"github.com/zaydek/retro/pkg/stdio_logger"
	"github.com/zaydek/retro/pkg/terminal"
	"github.com/zaydek/retro/pkg/vs"
	"github.com/zaydek/retro/pkg/watch"
)

const (
	WWW_DIR = "www"
	SRC_DIR = "src"
	OUT_DIR = "out"
)

var EPOCH = time.Now()

var cyan = func(str string) string { return pretty.Accent(str, terminal.Cyan) }
var magenta = func(str string) string { return pretty.Accent(str, terminal.Magenta) }

////////////////////////////////////////////////////////////////////////////////

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
			stdin <- ipc.Request{Kind: "rebuild"}
		}
	}()

	var once sync.Once
	go func() {
		stdin <- ipc.Request{Kind: "dev"}
		for {
			select {
			case out := <-stdout:
				once.Do(func() { ready <- struct{}{} })
				var res BuildResponse
				if err := json.Unmarshal(out.Data, &res); err != nil {
					panic(err)
				}
				dev <- res
			case err := <-stderr:
				stdio_logger.Stderr(err)
				// transformed := stdio_logger.TransformStderr(err)
				// fmt.Println(string(terminal_to_html.Render([]byte(transformed))))
				// os.Exit(1)
			}
		}
	}()

	r.Serve(ServerOptions{
		Stdin:     stdin,
		DevEvents: dev,
		Ready:     ready,
	})
}

////////////////////////////////////////////////////////////////////////////////

func (r Runner) Build() {
	fmt.Println("TODO")
}

////////////////////////////////////////////////////////////////////////////////

type ServerOptions struct {
	Stdin     chan ipc.Request
	DevEvents chan BuildResponse
	Ready     chan struct{}
}

func (r Runner) Serve(opt ServerOptions) {
	var res BuildResponse

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Server error (500)
		if res.Dirty() {
			fmt.Fprintln(w, res.HTML())
			return
		}
		// opt.Stdin <- ipc.Request{Kind: "rebuild"}
		path := getFSPath(r.URL.Path)
		if ext := filepath.Ext(path); ext == ".html" {
			http.ServeFile(w, r, filepath.Join(OUT_DIR, "index.html"))
			return
		}
		http.ServeFile(w, r, filepath.Join(OUT_DIR, path))
	})

	if opt.DevEvents != nil {
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
				case res = <-opt.DevEvents:
					bstr, err := json.Marshal(res)
					if err != nil {
						panic(err)
					}
					fmt.Fprintf(w, "event: reload\ndata: %s\n\n", string(bstr))
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

	dur := terminal.Dimf("(%s)", pretty.Duration(time.Since(EPOCH)))
	stdio_logger.Stdout(cyan(fmt.Sprintf("Ready on port '%d' %s", r.getPort(), dur)))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", r.getPort()), nil); err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////////////////////

var pkg Package

type Package struct {
	react                string
	react_dom            string
	retro                string
	retro_store          string
	retro_browser_router string
}

//go:embed pkg.txt
var contents string

func Run() {
	// Parse 'pkg.txt' and set pkg
	pkgMap, err := vs.Parse(contents)
	if err != nil {
		panic(err)
	}

	pkg = Package{
		react:                pkgMap["react"],
		react_dom:            pkgMap["react-dom"],
		retro:                pkgMap["@zaydek/retro"],
		retro_store:          pkgMap["@zaydek/retro-store"],
		retro_browser_router: pkgMap["@zaydek/retro-browser-router"],
	}

	cmd, err := cli.ParseCLIArguments()
	switch err {
	case cli.VersionError:
		fmt.Fprintln(os.Stdout, pkg.retro)
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
