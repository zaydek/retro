package retro

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
		// Blocks the serve command
		ready = make(chan struct{})

		// Sends messages to the serve command. Note that the dev channel needs to
		// be buffered so send operations are non-blocking.
		dev = make(chan Message, 1)

		// Orchestrates `copyIndexHTMLEntryPoint`
		once sync.Once
	)

	stdin <- "build"

	go func() {
		// TOOD: Where do we put `stdin <- "done"`?
	loop:
		for {
			select {
			case line := <-stdout:
				var message Message
				if err := json.Unmarshal([]byte(line), &message); err != nil {
					// Log unmarshal errors so users can debug plugins, etc.
					fmt.Println(decorateStdoutLine(line))
				} else {
					once.Do(func() {
						entries := entryPoints{clientCSS: "client.css", vendorJS: "vendor.js", clientJS: "client.js"}
						if err := copyIndexHTMLEntryPoint(entries); err != nil {
							// Panic because of the goroutine
							panic(fmt.Errorf("copyIndexHTMLEntryPoint: %w", err))
						}
						ready <- struct{}{}
					})
					// stdin <- "done"
					dev <- message
					break loop
				}
			case text := <-stderr:
				fmt.Fprintln(os.Stderr, decorateStderrText(text))
				stdin <- "done"
				os.Exit(1)
			}
		}
	}()

	// // DEBUG
	// bstr, err := json.MarshalIndent(message, "", "  ")
	// if err != nil {
	// 	return fmt.Errorf("json.MarshalIndent: %w", err)
	// }
	// fmt.Println(string(bstr))

	go func() {
		for result := range watch.Directory(RETRO_SRC_DIR, 100*time.Millisecond) {
			if result.Err != nil {
				// Panic because of the goroutine
				panic(fmt.Errorf("watch.Directory: %w", result.Err))
			}
			stdin <- "rebuild"
		}
	}()

	<-ready
	if err := a.Serve(ServeOptions{WarmUpFlag: false, Dev: dev}); err != nil {
		return fmt.Errorf("a.Serve: %w", err)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

type ServeOptions struct {
	WarmUpFlag bool
	Dev        chan Message
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
	bstr, err := os.ReadFile(filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR, "index.html"))
	if err != nil {
		return err
	}
	// Add the server sent events (SSE) stub
	contents := strings.Replace(
		string(bstr),
		"</body>",
		fmt.Sprintf("\t%s\n\t</body>", serverSentEventsStub),
		1,
	)
	fmt.Println(contents)

	var message Message
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/__dev__" {
			return
		}

		// // DEBUG
		// bstr, err := json.MarshalIndent(done, "", "  ")
		// if err != nil {
		// 	panic(fmt.Errorf("json.MarshalIndent: %w", err))
		// }
		// fmt.Println(string(bstr))

		// 500 Server error (esbuild errors)
		if message.Data.Vendor.IsDirty() {
			// Mirror errors to the browser and the terminal
			terminal.Clear(os.Stderr)
			fmt.Fprint(w, message.Data.Vendor.HTML())
			fmt.Fprint(os.Stderr, message)
			return
		} else if message.Data.Client.IsDirty() {
			// Mirror errors to the browser and the terminal
			terminal.Clear(os.Stderr)
			fmt.Fprint(w, message.Data.Client.HTML())
			fmt.Fprint(os.Stderr, message)
			return
		}

		url := getFilesystemPath(r.URL.Path)
		if extension := filepath.Ext(url); extension != "" && extension != ".html" {
			// 200 OK - Serve non-HTML
			http.ServeFile(w, r, filepath.Join(RETRO_OUT_DIR, url))
			return
		} else if a.getCommandKind() == KindDevCommand {
			// 200 OK - Serve HTML + server-sent events (SSE)
			fmt.Fprint(w, contents)
			if err := buildSuccess(a.getPort()); err != nil {
				panic(fmt.Errorf("buildSuccess: %w", err))
			}
		} else {
			// 200 OK - Serve HTML
			http.ServeFile(w, r, filepath.Join(RETRO_OUT_DIR, "index.html"))
			if err := buildSuccess(a.getPort()); err != nil {
				panic(fmt.Errorf("buildSuccess: %w", err))
			}
		}
	})

	// Add handler for dev events
	if a.getCommandKind() == KindDevCommand {
		http.HandleFunc("/__dev__", func(w http.ResponseWriter, r *http.Request) {
			// Add headers for server-sent events (SSE)
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			flusher, ok := w.(http.Flusher)
			if !ok {
				panic("w.(http.Flusher)")
			}
			// fmt.Println("Here")
			for {
				select {
				case message = <-options.Dev:
					// fmt.Println("A reload event occurred")
					fmt.Fprint(w, "event: reload\ndata\n\n")
					flusher.Flush()
				case <-r.Context().Done():
					return
				}
			}
		})
	}

	go func() {
		time.Sleep(10 * time.Millisecond)
		buildSuccess(a.getPort())
	}()

	for {
		err := http.ListenAndServe(fmt.Sprintf(":%d", a.getPort()), nil)
		if err != nil {
			if err.Error() == fmt.Sprintf("listen tcp :%d: bind: address already in use", a.getPort()) {
				a.setPort(a.getPort() + 1)
				continue
			}
			panic(err)
		}
		break
	}

	return nil
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
	case *cli.DevCommand:
		if err := app.Dev(DevOptions{WarmUpFlag: true}); err != nil {
			panic(fmt.Errorf("app.Dev: %w", err))
		}
		// case cli.BuildCommand:
		// 	app.Build(BuildOptions{WarmUpFlag: true})
		// case cli.ServeCommand:
		// 	app.Serve(ServeOptions{})
	}
}
