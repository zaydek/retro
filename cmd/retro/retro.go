package retro

import (
	_ "embed"
	"io/ioutil"
	"log"
	"net"
	"sort"
	"strings"

	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/zaydek/retro/cmd/format"
	"github.com/zaydek/retro/cmd/retro/cli"
	"github.com/zaydek/retro/pkg/ipc"
	"github.com/zaydek/retro/pkg/terminal"
	"github.com/zaydek/retro/pkg/watch"
)

var EPOCH = time.Now()

var cyan = func(str string) string { return format.Accent(str, terminal.Cyan) }

// getBase gets the executable directory name (basename). This directory
// changes depending on development or production.
func getBase() (string, error) {
	exec, err := os.Executable()
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(filepath.Base(exec), "main") {
		return os.Getwd()
	}
	// Get node_modules/.bin/@zaydek/bin/retro not node_modules/.bin/retro.
	return filepath.EvalSymlinks(exec)
}

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
			fmt.Fprintln(os.Stderr, format.Error(err.Error()))
			os.Exit(1)
		default:
			if err != nil {
				panic(err)
			}
		}
	}

	root, err := getBase()
	if err != nil {
		panic(err)
	}

	stdin, stdout, stderr, err := ipc.NewCommand("node", filepath.Join(filepath.Dir(root), "scripts/backend.esbuild.js"))
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
					if err := copyHTMLEntryPoint("vendor.js", "bundle.js", "bundle.css"); err != nil {
						panic(err)
					}
					ready <- struct{}{}
				})
				dev <- res
			case err := <-stderr:
				fmt.Fprint(os.Stderr, err)
			}
		}
	}()

	r.Serve(ServeOptions{Dev: dev, Ready: ready})
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
			fmt.Fprintln(os.Stderr, format.Error(err.Error()))
			os.Exit(1)
		default:
			if err != nil {
				panic(err)
			}
		}
	}

	root, err := getBase()
	if err != nil {
		panic(err)
	}

	stdin, stdout, stderr, err := ipc.NewCommand("node", filepath.Join(filepath.Dir(root), "scripts/backend.esbuild.js"))
	if err != nil {
		panic(err)
	}

	stdin <- ipc.Request{Kind: "build"}

	select {
	case out := <-stdout:
		// FIXME: stdout messages e.g. `console.log` from retro.config.js should not
		// be treated as errors if they fail to unmarshal. The problem is that
		// ipc.Message needs to be more blunt and simply provide a plaintext
		// interface for interacting with stdout and stderr.
		//
		// See https://github.com/zaydek/retro/issues/8.
		var res BackendResponse
		if err := json.Unmarshal(out.Data, &res); err != nil {
			panic(err)
		}
		if res.Dirty() {
			fmt.Fprint(os.Stderr, res)
			os.Exit(1)
		}
		vendorDotJS, bundleDotJS, bundleDotCSS := res.getChunkedNames()
		if err := copyHTMLEntryPoint(vendorDotJS, bundleDotJS, bundleDotCSS); err != nil {
			panic(err)
		}
	case err := <-stderr:
		fmt.Fprint(os.Stderr, err)
	}

	infos, err := ls(OUT_DIR)
	if err != nil {
		panic(err)
	}
	sort.Sort(infos)

	var sum int64
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
			strings.Repeat(" ", 40-len(v.path)),
			terminal.Dimf("(%s)", byteCount(v.size)),
		)

		if !strings.HasSuffix(v.path, ".map") {
			sum += v.size
		}
	}

	fmt.Println(strings.Repeat(" ", 40) + terminal.Dimf("(%s sum)", byteCount(sum)))
	fmt.Println()
	fmt.Println(terminal.Dimf("(%s)", time.Since(EPOCH)))
}

////////////////////////////////////////////////////////////////////////////////

type ServeOptions struct {
	Preflight bool

	Dev   chan BackendResponse
	Ready chan struct{}
}

// https://stackoverflow.com/a/37382208
func getIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func buildSuccess(port int) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var (
		base = filepath.Base(cwd)
		ip   = getIP()
	)

	terminal.Clear(os.Stdout)
	fmt.Println(terminal.Green("Compiled successfully!") + `

You can now view ` + terminal.Bold(base) + ` in the browser.

  ` + terminal.Bold("Local:") + `            ` + fmt.Sprintf("http://localhost:%s", terminal.Bold(port)) + `
  ` + terminal.Bold("On Your Network:") + `  ` + fmt.Sprintf("http://%s:%s", ip, terminal.Bold(port)) + `

Note that the development build is not optimized.
To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.
` /* EOF */)
}

func (r Runner) Serve(opt ServeOptions) {
	if opt.Preflight {
		_, err := r.preflight()
		switch err.(type) {
		case HTMLError:
			fmt.Fprintln(os.Stderr, format.Error(err.Error()))
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
		contents = strings.Replace(contents, "</body>", fmt.Sprintf("\t%s\n\t</body>", devStub), 1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/~dev" {
			return
		}

		// 500 Server error
		if res.Dirty() {
			terminal.Clear(os.Stderr)
			fmt.Fprint(w, res.HTML())
			fmt.Fprint(os.Stderr, res)
			return
		}
		// 200 OK - Serve non-index.html
		path := getFSPath(req.URL.Path)
		if ext := filepath.Ext(path); ext != "" && ext != ".html" {
			http.ServeFile(w, req, filepath.Join(OUT_DIR, path))
			return
		}
		// 200 OK - Serve index.html
		if r.getCommandKind() == KindDevCommand {
			fmt.Fprint(w, contents)
			buildSuccess(r.getPort())
		} else {
			http.ServeFile(w, req, filepath.Join(OUT_DIR, "index.html"))
			buildSuccess(r.getPort())
		}
	})

	if r.getCommandKind() != KindServeCommand {
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

	var (
		port    = r.getPort()
		getPort = func() int { return port }
	)

	go func() {
		time.Sleep(10 * time.Millisecond)
		buildSuccess(getPort())
	}()

	for {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			if err.Error() == fmt.Sprintf("listen tcp :%d: bind: address already in use", port) {
				port++
				continue
			}
			panic(err)
		}
		break
	}
}

////////////////////////////////////////////////////////////////////////////////

func Run() {
	cmd, err := cli.ParseCLIArguments()
	switch err {
	case cli.ErrVersion:
		fmt.Println(os.Getenv("RETRO_VERSION"))
		return
	case cli.ErrUsage:
		fallthrough
	case cli.ErrHelp:
		fmt.Println(format.Pad(format.Tabs(cyan(usage))))
		return
	}

	switch err.(type) {
	case cli.CommandError:
		fmt.Fprintln(os.Stderr, format.Error(err.Error()))
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
		run.Serve(ServeOptions{})
	}
}
