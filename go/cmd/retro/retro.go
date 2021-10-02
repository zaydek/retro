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

func (a *App) Dev(options DevOptions) error {
	if options.WarmUpFlag {
		var entryPointErr EntryPointError
		if err := warmUp(a.getCommandKind()); err != nil {
			if errors.As(err, &entryPointErr) {
				fmt.Fprintln(os.Stderr, format.Error(err))
				os.Exit(1)
			} else {
				return fmt.Errorf("warmUp: %w", err)
			}
		}
	}

	backendJS := filepath.Join(__dirname, "scripts/backend.esbuild.js")
	stdin, stdout, stderr, err := ipc.NewCommand("node", backendJS)
	if err != nil {
		return fmt.Errorf("ipc.NewCommand: %w", err)
	}

	var (
		// Blocks the serve command
		ready = make(chan struct{})

		// Sends messages to `/__dev__` HTTP handler
		dev = make(chan Message)
	)

	go func() {
		stdin <- "build"

		var once sync.Once
		for {
			select {
			case line := <-stdout:
				var message Message
				if err := json.Unmarshal([]byte(line), &message); err != nil {
					// Log unmarshal errs so users can debug plugins, etc.
					fmt.Println(formatStdoutLine(line))
				} else {
					once.Do(func() {
						entries := message.getChunkedEntrypoints()
						if err := copyIndexHTMLEntryPoint(entries); err != nil {
							panic(fmt.Errorf("copyIndexHTMLEntryPoint: %w", err))
						}
						ready <- struct{}{}
					})
					dev <- message
				}
			case text := <-stderr:
				fmt.Fprintln(os.Stderr, formatStderrText(text))
				os.Exit(1)
			}
		}
	}()

	go func() {
		for result := range watch.Directory(RETRO_SRC_DIR, 100*time.Millisecond) {
			if result.Err != nil {
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

type BuildOptions struct {
	WarmUpFlag bool
}

func (a *App) Build(options BuildOptions) error {
	if options.WarmUpFlag {
		var entryPointErr EntryPointError
		if err := warmUp(a.getCommandKind()); err != nil {
			if errors.As(err, &entryPointErr) {
				fmt.Fprintln(os.Stderr, format.Error(err))
				os.Exit(1)
			} else {
				return fmt.Errorf("warmUp: %w", err)
			}
		}
	}

	backendJS := filepath.Join(__dirname, "scripts/backend.esbuild.js")
	stdin, stdout, stderr, err := ipc.NewCommand("node", backendJS)
	if err != nil {
		return fmt.Errorf("ipc.NewCommand: %w", err)
	}

	stdin <- "build"

loop:
	for {
		select {
		case line := <-stdout:
			var message Message
			if err := json.Unmarshal([]byte(line), &message); err != nil {
				// Log unmarshal errs so users can debug plugins, etc.
				fmt.Println(formatStdoutLine(line))
			} else {
				if bundle := message.GetDirty(); bundle.IsDirty() {
					fmt.Print(bundle.String())
					os.Exit(1)
				}
				entries := message.getChunkedEntrypoints()
				if err := copyIndexHTMLEntryPoint(entries); err != nil {
					return fmt.Errorf("copyIndexHTMLEntryPoint: %w", err)
				}
				break loop
			}
		case text := <-stderr:
			fmt.Fprintln(os.Stderr, formatStderrText(text))
			os.Exit(1)
		}
	}

	str, err := buildBuildCommandSuccess(RETRO_OUT_DIR)
	if err != nil {
		return fmt.Errorf("buildLog: %w", err)
	}
	fmt.Print(str)

	return nil
}

////////////////////////////////////////////////////////////////////////////////

type ServeOptions struct {
	WarmUpFlag bool
	Dev        chan Message
}

func (a *App) Serve(options ServeOptions) error {
	if options.WarmUpFlag {
		var entryPointErr EntryPointError
		if err := warmUp(a.getCommandKind()); err != nil {
			if errors.As(err, &entryPointErr) {
				fmt.Fprintln(os.Stderr, format.Error(err))
				os.Exit(1)
			} else {
				return fmt.Errorf("warmUp: %w", err)
			}
		}
	}

	// out/index.html
	bstr, err := os.ReadFile(filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR, "index.html"))
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}
	contents := strings.Replace(
		string(bstr),
		"</body>",
		// Add server-sent events (SSE)
		fmt.Sprintf("\t%s\n\t</body>", serverSentEventsStub),
		1,
	)

	var message Message

	// Path for HTML and non-HTML resources
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 500 Server error
		if bundle := message.GetDirty(); bundle.IsDirty() {
			terminal.Clear(os.Stdout)
			fmt.Print(bundle.HTML())
			fmt.Print(bundle.String())
		}
		// 200 OK - Non-HTML resources
		path := getFilesystemPath(r.URL.Path)
		if extension := filepath.Ext(path); extension != "" && extension != ".html" {
			http.ServeFile(w, r, filepath.Join(RETRO_OUT_DIR, path))
			return
		}
		// 200 OK - HTML resources
		if a.getCommandKind() == KindDevCommand {
			fmt.Fprint(w, contents)
		} else {
			http.ServeFile(w, r, filepath.Join(RETRO_OUT_DIR, "index.html"))
		}
		terminal.Clear(os.Stdout)
		fmt.Println(buildServeCommandSuccess(a.getPort()))
	})

	// Path for server-sent events (SSE)
	http.HandleFunc("/__dev__", func(w http.ResponseWriter, r *http.Request) {
		if a.getCommandKind() == KindDevCommand {
			// Add headers for server-sent events (SSE)
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			flusher, ok := w.(http.Flusher)
			if !ok {
				panic("w.(http.Flusher)")
			}
			for {
				select {
				case message = <-options.Dev:
					fmt.Fprint(w, "event: reload\ndata\n\n")
					flusher.Flush()
				case <-r.Context().Done():
					return
				}
			}
		}
	})

	var port = a.getPort()

	go func() {
		time.Sleep(10 * time.Millisecond)
		terminal.Clear(os.Stdout)
		fmt.Println(buildServeCommandSuccess(port))
	}()

	for {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			if err.Error() == fmt.Sprintf("listen tcp :%d: bind: address already in use", port) {
				port++
				continue
			} else {
				return fmt.Errorf("http.ListenAndServe: %w", err)
			}
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
		fmt.Fprintln(os.Stderr, format.Error(err))
		os.Exit(1)
	default:
		if err != nil {
			panic(fmt.Errorf("cli.ParseCLIArguments: %w", err))
		}
	}

	var (
		app         = &App{Command: command}
		commandName string
	)

	switch app.Command.(type) {
	case cli.DevCommand:
		err = app.Dev(DevOptions{WarmUpFlag: true})
		commandName = "app.Dev"
	case cli.BuildCommand:
		err = app.Build(BuildOptions{WarmUpFlag: true})
		commandName = "app.Build"
	case cli.ServeCommand:
		err = app.Serve(ServeOptions{})
		commandName = "app.Serve"
	}
	if err != nil {
		panic(fmt.Errorf("%s: %w", commandName, err))
	}
}
