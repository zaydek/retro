package retro

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zaydek/retro/pkg/ipc"
	"github.com/zaydek/retro/pkg/stdio_logger"
	"github.com/zaydek/retro/pkg/terminal"
)

const (
	WWW_DIR = "www"
	SRC_DIR = "src"
	OUT_DIR = "out"
)

func getBrowserPath(url string) string {
	out := url
	if strings.HasSuffix(url, "/index.html") {
		out = out[:len(out)-len("index.html")] // Keep "/"
	} else if strings.HasSuffix(url, "/index") {
		out = out[:len(out)-len("index")] // Keep "/"
	} else if strings.HasSuffix(url, ".html") {
		out = out[:len(out)-len(".html")]
	}
	return out
}

func getFSPath(url string) string {
	out := url
	if strings.HasSuffix(url, "/") {
		out += "index.html"
	} else if strings.HasSuffix(url, "/index") {
		out += ".html"
	} else if ext := filepath.Ext(url); ext == "" {
		out += ".html"
	}
	return out
}

func Start() {
	os.Setenv("WWW_DIR", "www")
	os.Setenv("SRC_DIR", "src")
	os.Setenv("OUT_DIR", "out")

	stdin, stdout, stderr, err := ipc.NewCommand("node", "scripts/bundle.js")
	if err != nil {
		panic(err)
	}

	// service := ipc.Service{
	// 	Stdin:  stdin,
	// 	Stdout: stdout,
	// 	Stderr: stderr,
	// }

	dev := make(chan BuildResponse, 1)

	go func() {
		stdin <- ipc.Request{Kind: "build"}
		for {
			select {
			case out := <-stdout:
				var res BuildResponse
				if err := json.Unmarshal(out.Data, &res); err != nil {
					panic(err)
				}
				dev <- res
			case err := <-stderr:
				if err != "" {
					stdio_logger.Stderr(err)
				}
			}
		}
	}()

	Serve(ServerOptions{DevEvents: dev})

	// ch := make(chan struct{})
	// <-ch
}

type ServerOptions struct {
	DevEvents chan BuildResponse
}

func Serve(opts ServerOptions) {
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

	var port = 8000
	if envPort := os.Getenv("PORT"); envPort != "" {
		port, _ = strconv.Atoi(envPort)
	}

	stdio_logger.Stdout(terminal.Boldf("Ready on port %s.", terminal.Cyanf("'%d'", port)))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
