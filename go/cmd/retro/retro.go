package retro

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/zaydek/retro/go/cmd/format"
	"github.com/zaydek/retro/go/cmd/retro/cli"
	"github.com/zaydek/retro/go/pkg/ipc"
	"github.com/zaydek/retro/go/pkg/terminal"
)

var EPOCH = time.Now()

var cyan = func(str string) string { return format.Accent(str, terminal.Cyan) }

////////////////////////////////////////////////////////////////////////////////

type DevOptions struct {
	Preflight bool
}

func fatalUserError(err error) {
	// TODO: Clean this up; this is too vague
	fmt.Fprintln(os.Stderr, format.Error(err.Error()))
	os.Exit(1)
}

func (a *App) Dev(options DevOptions) {
	if options.Preflight {
		switch err := warmUp(a.getCommandKind()); err.(type) {
		case EntryPointError:
			fatalUserError(err)
		default:
			if err != nil {
				panic(fmt.Errorf("warmUp: %w", err))
			}
		}
	}

	stdin, stdout, stderr, err := ipc.NewCommand("node", filepath.Join(__dirname, "node/backend.esbuild.js"))
	if err != nil {
		panic(fmt.Errorf("ipc.NewCommand: %w", err))
	}

	var (
		isReadyToServe = make(chan struct{})

		// TODO: Why is the dev channel buffered?
		dev = make(chan BackendResponse, 1)
	)

	// go func() {
	// 	// TODO: In theory this shouldn't fire until the user does something; we
	// 	// need to check this doesn't fire eagerly
	// 	for result := range watch.Directory(RETRO_SRC_DIR, 100*time.Millisecond) {
	// 		if result.Err != nil {
	// 			panic(fmt.Errorf("watch.Directory: %w", result.Err))
	// 		}
	// 		stdin <- "rebuild"
	// 	}
	// }()

	go func() {
		var (
			message BackendResponse

			// For `transformAndCopyIndexHTMLEntryPoint`; extracts the cache-friendly
			// filenames for `src/index.css`, `src/index.js`, and `src/App.js`
			once sync.Once
		)

		stdin <- "build"

		// Technically we don't need a for-loop here except that user plugins can
		// log to stdout or stderr repeatedly
		for {
			select {
			case line := <-stdout:
				if err := json.Unmarshal([]byte(line), &message); err != nil {
					// Log unmarshal errors as stdout so users can debug plugins, etc.
					fmt.Println(decorateStdoutLine(line))
					continue
				}
				once.Do(func() {
					// For development, there's no reason to cache-bust the vendor or
					// client bundles; pass the canonical filenames as-is
					if err := transformAndCopyIndexHTMLEntryPoint("vendor.js", "client.js", "client.css"); err != nil {
						panic(fmt.Errorf("transformAndCopyIndexHTMLEntryPoint: %w"))
					}
					isReadyToServe <- struct{}{}
				})
				dev <- message

				// Break the select statement
				stdin <- "done"
				return
			case text := <-stderr:
				fmt.Fprintln(os.Stderr, decorateStderrText(text))

				// Break the select statement
				stdin <- "done"
				return
			}
		}
	}()

	// a.Serve(ServeOptions{Dev: dev, Ready: ready})
}

////////////////////////////////////////////////////////////////////////////////

// type BuildOptions struct {
// 	Preflight bool
// }
//
// func (r *App) Build(opt BuildOptions) {
// 	var transformAndCopyIndexHTMLEntryPoint func(string, string, string) error
// 	if opt.Preflight {
// 		var err error
// 		transformAndCopyIndexHTMLEntryPoint, err = r.warmUp()
// 		switch err.(type) {
// 		case HTMLError:
// 			fmt.Fprintln(os.Stderr, format.Error(err.Error()))
// 			os.Exit(1)
// 		default:
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}
//
// 	root, err := getBase()
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	stdin, stdout, stderr, err := ipc.NewCommand("node", filepath.Join(filepath.Dir(root), "scripts/backend.esbuild.js"))
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	stdin <- ipc.Request{Kind: "build"}
//
// 	select {
// 	case out := <-stdout:
// 		// FIXME: stdout messages e.g. `console.log` from retro.config.js should not
// 		// be treated as errors if they fail to unmarshal. The problem is that
// 		// ipc.Message needs to be more blunt and simply provide a plaintext
// 		// interface for interacting with stdout and stderr.
// 		//
// 		// See https://github.com/zaydek/retro/issues/8.
// 		var res BackendResponse
// 		if err := json.Unmarshal(out.Data, &res); err != nil {
// 			panic(err)
// 		}
// 		if res.Dirty() {
// 			fmt.Fprint(os.Stderr, res)
// 			os.Exit(1)
// 		}
// 		vendorDotJS, bundleDotJS, bundleDotCSS := res.getChunkedNames()
// 		if err := transformAndCopyIndexHTMLEntryPoint(vendorDotJS, bundleDotJS, bundleDotCSS); err != nil {
// 			panic(err)
// 		}
// 	case err := <-stderr:
// 		fmt.Fprint(os.Stderr, err)
// 	}
//
// 	infos, err := ls(RETRO_OUT_DIR)
// 	if err != nil {
// 		panic(err)
// 	}
// 	sort.Sort(infos)
//
// 	var sum int64
// 	for _, v := range infos {
// 		var color = terminal.Normal
// 		if strings.HasSuffix(v.path, ".html") {
// 			color = terminal.Normal
// 		} else if strings.HasSuffix(v.path, ".js") || strings.HasSuffix(v.path, ".js.map") {
// 			color = terminal.Yellow
// 		} else if strings.HasSuffix(v.path, ".css") || strings.HasSuffix(v.path, ".css.map") {
// 			color = terminal.Cyan
// 		} else {
// 			color = terminal.Dim
// 		}
//
// 		fmt.Printf("%v%s%v\n",
// 			color(v.path),
// 			strings.Repeat(" ", 40-len(v.path)),
// 			terminal.Dimf("(%s)", byteCount(v.size)),
// 		)
//
// 		if !strings.HasSuffix(v.path, ".map") {
// 			sum += v.size
// 		}
// 	}
//
// 	fmt.Println(strings.Repeat(" ", 40) + terminal.Dimf("(%s sum)", byteCount(sum)))
// 	fmt.Println()
// 	fmt.Println(terminal.Dimf("(%s)", time.Since(EPOCH)))
// }

////////////////////////////////////////////////////////////////////////////////

// type ServeOptions struct {
// 	Preflight bool
//
// 	Dev   chan BackendResponse
// 	Ready chan struct{}
// }
//
// // https://stackoverflow.com/a/37382208
// //
// // FIXME: This doesn't support the use-case that the user isn't connected to the
// // internet, which makes Retro unusable for internet-less development
// func getIP() net.IP {
// 	conn, err := net.Dial("udp", "8.8.8.8:80")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Close()
// 	localAddr := conn.LocalAddr().(*net.UDPAddr)
// 	return localAddr.IP
// }
//
// func buildSuccess(port int) {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	var (
// 		base = filepath.Base(cwd)
// 		ip   = getIP()
// 	)
//
// 	terminal.Clear(os.Stdout)
// 	fmt.Println(terminal.Green("Compiled successfully!") + `
//
// You can now view ` + terminal.Bold(base) + ` in the browser.
//
//   ` + terminal.Bold("Local:") + `            ` + fmt.Sprintf("http://localhost:%s", terminal.Bold(port)) + `
//   ` + terminal.Bold("On Your Network:") + `  ` + fmt.Sprintf("http://%s:%s", ip, terminal.Bold(port)) + `
//
// Note that the development build is not optimized.
// To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.
// ` /* EOF */)
// }
//
// func (r *App) Serve(opt ServeOptions) {
// 	if opt.Preflight {
// 		_, err := r.warmUp()
// 		switch err.(type) {
// 		case HTMLError:
// 			fmt.Fprintln(os.Stderr, format.Error(err.Error()))
// 			os.Exit(1)
// 		default:
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}
//
// 	if opt.Ready != nil {
// 		<-opt.Ready
// 	}
//
// 	var res BackendResponse
//
// 	// Add the dev stub
// 	var contents string
// 	if os.Getenv("ENV") == "development" {
// 		bstr, err := os.ReadFile(filepath.Join(RETRO_OUT_DIR, "index.html"))
// 		if err != nil {
// 			panic(err)
// 		}
// 		contents = string(bstr)
// 		contents = strings.Replace(contents, "</body>", fmt.Sprintf("\t%s\n\t</body>", serverSentEventsStub), 1)
// 	}
//
// 	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
// 		if req.URL.Path == "/~dev" {
// 			return
// 		}
//
// 		// 500 Server error
// 		if res.Dirty() {
// 			terminal.Clear(os.Stderr)
// 			fmt.Fprint(w, res.HTML())
// 			fmt.Fprint(os.Stderr, res)
// 			return
// 		}
// 		// 200 OK - Serve non-index.html
// 		path := getFilesystemPath(req.URL.Path)
// 		if ext := filepath.Ext(path); ext != "" && ext != ".html" {
// 			http.ServeFile(w, req, filepath.Join(RETRO_OUT_DIR, path))
// 			return
// 		}
// 		// 200 OK - Serve index.html
// 		if r.getCommandKind() == KindDevCommand {
// 			fmt.Fprint(w, contents)
// 			buildSuccess(r.getPort())
// 		} else {
// 			http.ServeFile(w, req, filepath.Join(RETRO_OUT_DIR, "index.html"))
// 			buildSuccess(r.getPort())
// 		}
// 	})
//
// 	if r.getCommandKind() != KindServeCommand {
// 		http.HandleFunc("/~dev", func(w http.ResponseWriter, req *http.Request) {
// 			w.Header().Set("Content-Type", "text/event-stream")
// 			w.Header().Set("Cache-Control", "no-cache")
// 			w.Header().Set("Connection", "keep-alive")
// 			flusher, ok := w.(http.Flusher)
// 			if !ok {
// 				panic("Internal error")
// 			}
// 			for {
// 				select {
// 				case res = <-opt.Dev:
// 					fmt.Fprint(w, "event: reload\ndata\n\n")
// 					flusher.Flush()
// 				case <-req.Context().Done():
// 					return
// 				}
// 			}
// 		})
// 	}
//
// 	var (
// 		port    = r.getPort()
// 		getPort = func() int { return port }
// 	)
//
// 	go func() {
// 		time.Sleep(10 * time.Millisecond)
// 		buildSuccess(getPort())
// 	}()
//
// 	for {
// 		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
// 		if err != nil {
// 			if err.Error() == fmt.Sprintf("listen tcp :%d: bind: address already in use", port) {
// 				port++
// 				continue
// 			}
// 			panic(err)
// 		}
// 		break
// 	}
// }

////////////////////////////////////////////////////////////////////////////////

var __dirname string

func Run() {
	var err error
	__dirname, err = getDirname()
	if err != nil {
		panic(fmt.Errorf("getDirname: %w", err))
	}

	// Parse the CLI arguments and guard sentinel errors
	command, err := cli.ParseCLIArguments()
	switch err {
	case cli.ErrVersion:
		fmt.Println(os.Getenv("RETRO_VERSION"))
		return
	case cli.ErrUsage:
		fallthrough
	case cli.ErrHelp:
		// TODO: Clean this up; this is too vague
		fmt.Println(format.Pad(format.Tabs(cyan(usage))))
		return
	}

	switch err.(type) {
	case cli.CommandError:
		// TODO: Clean this up; this is too vague
		fmt.Fprintln(os.Stderr, format.Error(err.Error()))
		os.Exit(1)
	default:
		if err != nil {
			panic(err)
		}
	}

	app := &App{Command: command}
	switch app.Command.(type) {
	case cli.DevCommand:
		app.Dev(DevOptions{Preflight: true})
		// case cli.BuildCommand:
		// 	app.Build(BuildOptions{Preflight: true})
		// case cli.ServeCommand:
		// 	app.Serve(ServeOptions{})
	}
}
