package retro

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/zaydek/retro/go/cmd/format"
	"github.com/zaydek/retro/go/cmd/retro/cli"
	"github.com/zaydek/retro/go/pkg/ipc"
	"github.com/zaydek/retro/go/pkg/terminal"
	"github.com/zaydek/retro/go/pkg/watch"
)

var EPOCH = time.Now()

var cyan = func(str string) string { return format.Accent(str, terminal.Cyan) }

////////////////////////////////////////////////////////////////////////////////

type DevOptions struct {
	WarmUpFlag bool
}

func fatalUserError(err error) {
	// TODO: Clean this up; this is too vague
	fmt.Fprintln(os.Stderr, format.Error(err.Error()))
	os.Exit(1)
}

func (a *App) Dev(options DevOptions) error {
	if options.WarmUpFlag {
		entryPointErrPointer := &EntryPointError{}
		if err := warmUp(a.getCommandKind()); err != nil {
			if errors.As(err, entryPointErrPointer) {
				fatalUserError(entryPointErrPointer)
			} else {
				return fmt.Errorf("warmUp: %w", err)
			}
		}
	}

	// Run the backend Node.js command
	stdin, stdout, stderr, err := ipc.NewCommand("node", filepath.Join(__dirname, "node/backend.esbuild.js"))
	if err != nil {
		return fmt.Errorf("ipc.NewCommand: %w", err)
	}

	var (
		done DoneMessage
		dev  = make(chan DoneMessage)

		// Orchestrates `copyIndexHTMLEntryPoint`
		once sync.Once
	)

	stdin <- "build"

	// Use a for-loop so plugins can log repeatedly
loop:
	for {
		select {
		case line := <-stdout:
			if err := json.Unmarshal([]byte(line), &done); err != nil {
				// Log unmarshal errors so users can debug plugins, etc.
				fmt.Println(decorateStdoutLine(line))
			} else {
				once.Do(func() {
					entries := entryPoints{clientCSS: "client.css", vendorJS: "vendor.js", clientJS: "client.js"}
					if err := copyIndexHTMLEntryPoint(entries); err != nil {
						// Panic because of the goroutine
						panic(fmt.Errorf("copyIndexHTMLEntryPoint: %w", err))
					}
				})
				// Done: Stop the Node.js runtime
				stdin <- "done"
				break loop
			}
		case text := <-stderr:
			fmt.Fprintln(os.Stderr, decorateStderrText(text))
			// Done: Stop the Node.js and Go runtime
			stdin <- "done"
			os.Exit(1)
		}
	}

	// DEBUG
	byteStr, err := json.MarshalIndent(done, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}
	fmt.Println(string(byteStr))

	go func() {
		for result := range watch.Directory(RETRO_SRC_DIR, 100*time.Millisecond) {
			if result.Err != nil {
				// Panic because of the goroutine
				panic(fmt.Errorf("watch.Directory: %w", result.Err))
			}
			stdin <- "rebuild"
		}
	}()

	if err := a.Serve(ServeOptions{WarmUpFlag: false, Dev: dev}); err != nil {
		return fmt.Errorf("a.Serve: %w", err)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

// type BuildOptions struct {
// 	WarmUpFlag bool
// }
//
// func (r *App) Build(opt BuildOptions) {
// 	var copyIndexHTMLEntryPoint func(string, string, string) error
// 	if opt.WarmUpFlag {
// 		var err error
// 		copyIndexHTMLEntryPoint, err = r.warmUp()
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
// 		if err := copyIndexHTMLEntryPoint(vendorDotJS, bundleDotJS, bundleDotCSS); err != nil {
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

type ServeOptions struct {
	WarmUpFlag bool
	Dev        chan DoneMessage
}

func (a *App) Serve(options ServeOptions) error {
	if options.WarmUpFlag {
		entryPointErrPointer := &EntryPointError{}
		if err := warmUp(a.getCommandKind()); err != nil {
			if errors.As(err, entryPointErrPointer) {
				fatalUserError(entryPointErrPointer)
			} else {
				return fmt.Errorf("warmUp: %w", err)
			}
		}
	}

	// www/index.html
	byteStr, err := os.ReadFile(filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR, "index.html"))
	if err != nil {
		return err
	}
	// Add the server sent events (SSE) stub
	contents := strings.Replace(
		string(byteStr),
		"</body>",
		fmt.Sprintf("\t%s\n\t</body>", serverSentEventsStub),
		1,
	)

	fmt.Print(contents)
	return nil

	//	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	//		if req.URL.Path == "/__dev__" {
	//			return
	//		}
	//
	//		// 500 Server error
	//		if done.Data.Vendor.IsDirty() || done.Data.Client.IsDirty() {
	//			terminal.Clear(os.Stderr) // TODO: Do we really want to clear the terminal?
	//			fmt.Fprint(w, done.HTML())
	//			fmt.Fprint(os.Stderr, done)
	//			return
	//		}
	//		// 200 OK - Serve non-index.html
	//		path := getFilesystemPath(req.URL.Path)
	//		if ext := filepath.Ext(path); ext != "" && ext != ".html" {
	//			http.ServeFile(w, req, filepath.Join(RETRO_OUT_DIR, path))
	//			return
	//		}
	//		// 200 OK - Serve index.html
	//		if a.getCommandKind() == KindDevCommand {
	//			fmt.Fprint(w, contents)
	//			buildSuccess(a.getPort())
	//		} else {
	//			http.ServeFile(w, req, filepath.Join(RETRO_OUT_DIR, "index.html"))
	//			buildSuccess(a.getPort())
	//		}
	//	})
	//
	//	if a.getCommandKind() != KindServeCommand {
	//		http.HandleFunc("/__dev__", func(w http.ResponseWriter, req *http.Request) {
	//			w.Header().Set("Content-Type", "text/event-stream")
	//			w.Header().Set("Cache-Control", "no-cache")
	//			w.Header().Set("Connection", "keep-alive")
	//			flusher, ok := w.(http.Flusher)
	//			if !ok {
	//				panic("Internal error")
	//			}
	//			for {
	//				select {
	//				case done = <-options.Dev:
	//					fmt.Fprint(w, "event: reload\ndata\n\n")
	//					flusher.Flush()
	//				case <-req.Context().Done():
	//					return
	//				}
	//			}
	//		})
	//	}
	//
	//	var (
	//		port    = a.getPort()
	//		getPort = func() int { return port }
	//	)
	//
	//	go func() {
	//		time.Sleep(10 * time.Millisecond)
	//		buildSuccess(getPort())
	//	}()
	//
	//	for {
	//		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	//		if err != nil {
	//			if err.Error() == fmt.Sprintf("listen tcp :%d: bind: address already in use", port) {
	//				port++
	//				continue
	//			}
	//			panic(err)
	//		}
	//		break
	//	}
}

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
		if err := app.Dev(DevOptions{WarmUpFlag: true}); err != nil {
			panic(fmt.Errorf("app.Dev: %w", err))
		}
		// case cli.BuildCommand:
		// 	app.Build(BuildOptions{WarmUpFlag: true})
		// case cli.ServeCommand:
		// 	app.Serve(ServeOptions{})
	}
}
