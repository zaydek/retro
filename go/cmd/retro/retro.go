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

type TimedMessage struct {
	msg Message
	dur time.Duration
}

func (a *App) Dev(options DevOptions) error {
	if options.WarmUpFlag {
		var entryPointErr EntryPointError
		if err := warmUp(a.getCommandKind()); err != nil {
			if errors.As(err, &entryPointErr) {
				fmt.Fprintln(os.Stderr, format.Stderr(err))
				os.Exit(1)
			} else {
				return err
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	stdin, stdout, stderr, err := ipc.NewPersistentCommand(ctx, "node", filepath.Join(__dirname, "scripts/backend.esbuild.js"))
	if err != nil {
		cancel()
		return err
	}
	defer cancel()

	var (
		dev   = make(chan TimedMessage)
		ready = make(chan struct{})
	)

	var tm time.Time
	go func() {
		tm = time.Now()
		stdin <- "build"
		var once sync.Once
		for {
			select {
			case line := <-stdout:
				var msg Message
				must(json.Unmarshal([]byte(line), &msg))
				once.Do(func() {
					entries := msg.getChunkedEntrypoints()
					must(copyIndexHTMLEntryPoint(entries))
					ready <- struct{}{}
				})
				dev <- TimedMessage{
					msg: msg,
					dur: time.Since(tm),
				}
			case text := <-stderr:
				fmt.Fprintln(os.Stderr, format.StderrIPC(text))
				cancel()
				os.Exit(1)
			}
		}
	}()

	go func() {
		for result := range watch.Directory(RETRO_SRC_DIR, 100*time.Millisecond) {
			must(result.Err)
			tm = time.Now()
			stdin <- "rebuild"
		}
	}()

	<-ready
	if err := a.Serve(ServeOptions{WarmUpFlag: false, Dev: dev}); err != nil {
		return err
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
				fmt.Fprintln(os.Stderr, format.Stderr(err))
				os.Exit(1)
			} else {
				return err
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	stdin, stdout, stderr, err := ipc.NewPersistentCommand(ctx, "node", filepath.Join(__dirname, "scripts/backend.esbuild.js"))
	if err != nil {
		cancel()
		return err
	}

	stdin <- "build"
	defer cancel()

loop:
	for {
		select {
		case line := <-stdout:
			var message Message
			if err := json.Unmarshal([]byte(line), &message); err != nil {
				// // Log unmarshal errors so users can debug plugins, etc.
				// fmt.Println(formatStdoutLine(line))
				return err
			}
			// Log to stdout and crash
			if dirty := message.GetDirty(); dirty.IsDirty() {
				fmt.Print(dirty.String())
				cancel()
				os.Exit(1)
			}
			entries := message.getChunkedEntrypoints()
			if err := copyIndexHTMLEntryPoint(entries); err != nil {
				return err
			}
			break loop
		case text := <-stderr:
			fmt.Fprintln(os.Stderr, format.StderrIPC(text))
			cancel()
			os.Exit(1)
		}
	}

	str, err := buildBuildSuccessString(RETRO_OUT_DIR)
	if err != nil {
		return err
	}
	fmt.Print(str)

	return nil
}

////////////////////////////////////////////////////////////////////////////////

type ServeOptions struct {
	WarmUpFlag bool
	Dev        chan TimedMessage
}

func (a *App) Serve(options ServeOptions) error {
	if options.WarmUpFlag {
		if err := setEnvAndGlobalVariables(KindServeCommand); err != nil {
			return err
		}
	}

	// out/index.html
	var contents string
	if a.getCommandKind() == KindDevCommand {
		bstr, err := os.ReadFile(filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR, "index.html"))
		if err != nil {
			return err
		}
		contents = strings.Replace(string(bstr), "</body>", fmt.Sprintf("\t%s\n\t</body>", htmlServerSentEvents), 1)
	}

	var (
		timedMsg = <-options.Dev
		logMsg   string
	)

	// Path for HTML and non-HTML resources
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log to stdout
		var nextLogMsg string
		if dirty := timedMsg.msg.GetDirty(); dirty.IsDirty() {
			nextLogMsg = dirty.String()
		} else {
			nextLogMsg = buildServeSuccessString(a.getPort(), timedMsg.dur)
		}
		if logMsg != nextLogMsg {
			logMsg = nextLogMsg
			terminal.Clear(os.Stdout)
			fmt.Println(logMsg)
		}
		// Log to the browser and eagerly return
		if dirty := timedMsg.msg.GetDirty(); dirty.IsDirty() {
			fmt.Fprintln(w, dirty.HTML())
			return
		}
		// Serve non-HTML
		path := getFilesystemPath(r.URL.Path)
		if ext := filepath.Ext(path); ext != "" && ext != ".html" {
			http.ServeFile(w, r, filepath.Join(RETRO_OUT_DIR, path))
			return
		}
		// Serve HTML
		if a.getCommandKind() == KindDevCommand {
			fmt.Fprint(w, contents)
			return
		}
		http.ServeFile(w, r, filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR, "index.html"))
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
				case timedMsg = <-options.Dev:
					fmt.Fprint(w, "event: reload\ndata\n\n")
					flusher.Flush()
				case <-r.Context().Done():
					return
				}
			}
		})
	}

	// Log to stdout
	var nextLogMsg string
	if dirty := timedMsg.msg.GetDirty(); dirty.IsDirty() {
		nextLogMsg = dirty.String()
	} else {
		nextLogMsg = buildServeSuccessString(a.getPort(), timedMsg.dur)
	}
	if logMsg != nextLogMsg {
		logMsg = nextLogMsg
		terminal.Clear(os.Stdout)
		fmt.Println(logMsg)
	}

	port := a.getPort()
	for {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			if err.Error() == fmt.Sprintf("listen tcp :%d: bind: address already in use", port) {
				port++
				continue
			} else {
				return err
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
	must(err)

	// Non-command errors
	command, err := cli.ParseCLIArguments()
	switch err {
	case cli.ErrVersion:
		fmt.Println(os.Getenv("RETRO_V_VERSION"))
		os.Exit(0)
	case cli.ErrUsage:
		fallthrough
	case cli.ErrHelp:
		fmt.Println(format.Stdout(usage))
		os.Exit(0)
	}

	// Command errors
	switch err.(type) {
	case cli.CommandError:
		fmt.Fprintln(os.Stderr, format.Stderr(err))
		os.Exit(1)
	default:
		must(err)
	}

	app := &App{Command: command}
	switch app.Command.(type) {
	case cli.DevCommand:
		err = app.Dev(DevOptions{WarmUpFlag: true})
	case cli.BuildCommand:
		err = app.Build(BuildOptions{WarmUpFlag: true})
	case cli.ServeCommand:
		err = app.Serve(ServeOptions{WarmUpFlag: true})
	}
	must(err)
}
