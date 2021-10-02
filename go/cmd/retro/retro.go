package retro

import (
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())
	stdin, stdout, stderr, err := ipc.NewCommand(ctx, "node", filepath.Join(__dirname, "scripts/backend.esbuild.js"))
	if err != nil {
		return fmt.Errorf("ipc.NewCommand: %w", err)
	}

	var (
		// Blocks the serve command
		ready = make(chan struct{})

		// Sends messages to the `/__dev__` HTTP handler
		dev = make(chan Message)
	)

	go func() {
		stdin <- "build"
		defer cancel()

		var once sync.Once
		for {
			select {
			case line := <-stdout:
				var message Message
				if err := json.Unmarshal([]byte(line), &message); err != nil {
					// // Log unmarshal errs so users can debug plugins, etc.
					// fmt.Println(formatStdoutLine(line))
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
				cancel()
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

	ctx, cancel := context.WithCancel(context.Background())
	stdin, stdout, stderr, err := ipc.NewCommand(ctx, "node", filepath.Join(__dirname, "scripts/backend.esbuild.js"))
	if err != nil {
		return fmt.Errorf("ipc.NewCommand: %w", err)
	}

	stdin <- "build"
	defer cancel()

loop:
	for {
		select {
		case line := <-stdout:
			var message Message
			if err := json.Unmarshal([]byte(line), &message); err != nil {
				// // Log unmarshal errs so users can debug plugins, etc.
				// fmt.Println(formatStdoutLine(line))
			} else {
				if dirty := message.GetDirty(); dirty.IsDirty() {
					fmt.Print(dirty.String())
					cancel()
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
			cancel()
			os.Exit(1)
		}
	}

	str, err := makeBuildSuccess(RETRO_OUT_DIR)
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
		if err := setEnv(KindServeCommand); err != nil {
			return fmt.Errorf("setEnv: %w", err)
		}
	}

	// out/index.html
	var contents string
	if a.getCommandKind() == KindDevCommand {
		bstr, err := os.ReadFile(filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR, "index.html"))
		if err != nil {
			return fmt.Errorf("os.ReadFile: %w", err)
		}
		contents = strings.Replace(string(bstr), "</body>", fmt.Sprintf("\t%s\n\t</body>", serverSentEventsStub), 1)
	}

	var (
		// Get the build message
		message = <-options.Dev

		// The logMessage log message
		logMessage string
	)

	// Path for HTML and non-HTML resources
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Dedupe the next log message (error)
		next := makeServeSuccess(a.getPort())
		if dirty := message.GetDirty(); dirty.IsDirty() {
			next = dirty.String()
		}
		if next != logMessage {
			logMessage = next
			terminal.Clear(os.Stdout)
			fmt.Println(logMessage)
			if dirty := message.GetDirty(); dirty.IsDirty() {
				fmt.Fprintln(w, dirty.HTML())
				// Eagerly return
				return
			}
		}
		// Serve non-HTML resources
		path := getFilesystemPath(r.URL.Path)
		if extension := filepath.Ext(path); extension != "" && extension != ".html" {
			http.ServeFile(w, r, filepath.Join(RETRO_OUT_DIR, path))
			return
		}
		// Serve nTML resources
		if a.getCommandKind() == KindDevCommand {
			fmt.Fprint(w, contents)
		} else {
			fmt.Println(filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR, "index.html"))
			http.ServeFile(w, r, filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR, "index.html"))
		}
	})

	// Path for server-sent events (SSE)
	if a.getCommandKind() == KindDevCommand {
		http.HandleFunc("/__dev__", func(w http.ResponseWriter, r *http.Request) {
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
		})
	}

	next := makeServeSuccess(a.getPort())
	if dirty := message.GetDirty(); dirty.IsDirty() {
		next = dirty.String()
	}
	if next != logMessage {
		logMessage = next
		terminal.Clear(os.Stdout)
		fmt.Println(logMessage)
	}

	port := a.getPort()
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
	decorate(&err, "getDirname")
	check(err)

	// Parse the CLI arguments and guard errors
	command, err := cli.ParseCLIArguments()
	switch err {
	case cli.ErrVersion:
		fmt.Println(os.Getenv("RETRO_VERSION"))
		return
	case cli.ErrUsage:
		fallthrough
	case cli.ErrHelp:
		fmt.Println(format.NonError(usage))
		return
	}

	switch err.(type) {
	case cli.CommandError:
		fmt.Fprintln(os.Stderr, format.Error(err))
		os.Exit(1)
	default:
		decorate(&err, "cli.ParseCLIArguments")
		check(err)
	}

	app := &App{Command: command}
	switch app.Command.(type) {
	case cli.DevCommand:
		err := app.Dev(DevOptions{WarmUpFlag: true})
		decorate(&err, "app.Dev")
	case cli.BuildCommand:
		err := app.Build(BuildOptions{WarmUpFlag: true})
		decorate(&err, "app.Build")
	case cli.ServeCommand:
		err := app.Serve(ServeOptions{WarmUpFlag: true})
		decorate(&err, "app.Build")
	}
	check(err)
}
