package retro

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	terminal_to_html "github.com/buildkite/terminal-to-html/v3"
	"github.com/zaydek/retro/cmd/pretty"
	"github.com/zaydek/retro/cmd/retro/cli"
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

var cyan = func(str string) string { return pretty.Accent(str, terminal.Cyan) }
var magenta = func(str string) string { return pretty.Accent(str, terminal.Magenta) }

////////////////////////////////////////////////////////////////////////////////
// % retro dev

func (r Runner) Dev() {
	os.Setenv("WWW_DIR", WWW_DIR)
	os.Setenv("SRC_DIR", SRC_DIR)
	os.Setenv("OUT_DIR", OUT_DIR)

	stdin, stdout, stderr, err := ipc.NewCommand("node", "scripts/backend.js")
	if err != nil {
		panic(err)
	}

	dev := make(chan BuildResponse, 1)

	// Setup a second watcher for .scss, etc. Tracked by
	// https://github.com/evanw/esbuild/issues/808.
	go func() {
		for watchRes := range watch.Directory(SRC_DIR, 100*time.Millisecond) {
			if watchRes.Error != nil {
				panic(err)
			}
			// fmt.Println("something changed")
			stdin <- ipc.Request{Kind: "rebuild"}
		}
	}()

	go func() {
		stdin <- ipc.Request{Kind: "build"}
		for {
			select {
			case out := <-stdout:
				// fmt.Println("stdout")
				var res BuildResponse
				if err := json.Unmarshal(out.Data, &res); err != nil {
					panic(err)
				}
				dev <- res
			case err := <-stderr:
				transformed := stdio_logger.TransformStderr(err)
				fmt.Println(string(terminal_to_html.Render([]byte(transformed))))
				// fmt.Println("stderr")
				stdio_logger.Stderr(err)
			}
		}
	}()

	r.Serve(ServerOptions{DevEvents: dev})
}

////////////////////////////////////////////////////////////////////////////////
// % retro build

func (r Runner) Build() {
	fmt.Println("TODO")
}

////////////////////////////////////////////////////////////////////////////////
// % retro serve

type ServerOptions struct {
	DevEvents chan BuildResponse
}

func (r Runner) Serve(opts ServerOptions) {
	var res BuildResponse

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Server error (500)
		if res.Dirty() {
			// http.Error(w, "500 server error", http.StatusInternalServerError)
			fmt.Fprintln(w, res.HTML())
			return
		}

		path := getFSPath(r.URL.Path)
		if ext := filepath.Ext(path); ext == ".html" {
			http.ServeFile(w, r, filepath.Join(OUT_DIR, "index.html"))
			return
		}
		http.ServeFile(w, r, filepath.Join(OUT_DIR, path))
	})

	if opts.DevEvents != nil {
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
				case res = <-opts.DevEvents:
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

	dur := terminal.Dimf("(%s)", pretty.Duration(time.Since(EPOCH)))
	stdio_logger.Stdout(cyan(fmt.Sprintf("Ready on port '%d' %s", r.getPort(), dur)))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", r.getPort()), nil); err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////////////////////
// % retro

func Run() {
	// TODO
	os.Setenv("RETRO_VERSION", "0.0.0")

	cmd, err := cli.ParseCLIArguments()
	switch err {
	case cli.VersionError:
		fmt.Println(os.Getenv("RETRO_VERSION"))
		return
	case cli.UsageError:
		fmt.Println(pretty.Inset(pretty.Spaces(cyan(usage))))
		os.Exit(1)
		return
	}

	switch err.(type) {
	case cli.CommandError:
		fmt.Fprintln(os.Stderr, pretty.Error(err.Error()))
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
		guarderr := guards()
		switch guarderr.(type) {
		case HTMLError:
			fmt.Fprintln(os.Stderr, pretty.Error(guarderr.Error()))
			os.Exit(1)
		default:
			if guarderr != nil {
				panic(guarderr)
			}
		}
		run.Dev()
	case cli.BuildCommand:
		os.Setenv("NODE_ENV", "production")
		guardErr := guards()
		switch guardErr.(type) {
		case HTMLError:
			fmt.Fprintln(os.Stderr, pretty.Error(guardErr.Error()))
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

	// ...
}
