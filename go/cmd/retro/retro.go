package retro

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	render "github.com/buildkite/terminal-to-html/v3"
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

// TODO
type DevError interface {
	IsDirty() bool
	String() string
	HTML() string
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
		dev   = make(chan DevError)
		ready = make(chan struct{})
	)

	// var tm time.Time // TODO
	go func() {
		// tm = time.Now() // TODO
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
				dev <- msg
			case text := <-stderr:
				fmt.Fprintln(os.Stderr, format.StderrIPC(text))
				cancel()
				os.Exit(1)
			}
		}
	}()

	go func() {
		for event := range watch.Directory(RETRO_SRC_DIR, 100*time.Millisecond) {
			must(event.Err)
			// tm = time.Now() // TODO
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
	defer cancel()

	tm := time.Now()
	stdin <- "build"

loop:
	for {
		select {
		case line := <-stdout:
			var msg Message
			if err := json.Unmarshal([]byte(line), &msg); err != nil {
				return err
			}
			// Log to stdout and crash
			if msg.IsDirty() {
				fmt.Print(msg.String())
				cancel()
				os.Exit(1)
			}
			entries := msg.getChunkedEntrypoints()
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

	str, err := buildBuildSuccessString(RETRO_OUT_DIR, time.Since(tm))
	if err != nil {
		return err
	}
	fmt.Print(str)

	return nil
}

////////////////////////////////////////////////////////////////////////////////

type RuntimeError struct {
	stack string
}

func newRuntimeError(stack string) RuntimeError {
	return RuntimeError{stack: stack}
}

func (r *RuntimeError) Reset() {
	(*r).stack = ""
}

func (r RuntimeError) IsDirty() bool {
	return r.stack != ""
}

// TODO: Upgrade to use esbuild-style errors
func (r RuntimeError) String() string {
	return r.stack
}

func (r RuntimeError) HTML() string {
	renderStr := string(
		render.Render(
			[]byte(
				r.String(),
			),
		),
	)
	return `<!DOCTYPE html>
<html lang="en">
	<head>
		<title>Runtime Error</title>
		<style>
:root {
	-webkit-font-smoothing: antialiased; /* macOS */
	-moz-osx-font-smoothing: grayscale;  /* Firefox */

	color: #c7c7c7;
	background-color: #000000;
}

code {
	font: 18px / 1.45
		"Monaco",   /* macOS */
		"Consolas", /* Windows */
		monospace;
}

a { color: unset; text-decoration: unset; }
a:hover { text-decoration: underline; }

.term-fg1 {
	font-weight: bold;
	color: #feffff;
}

.term-fg30 { color: #000000; }
.term-fg31 { color: #c91b00; }
.term-fg32 { color: #00c200; }
.term-fg33 { color: #c7c400; }
.term-fg34 { color: #0225c7; }
.term-fg35 { color: #c930c7; }
.term-fg36 { color: #00c5c7; }
.term-fg37 { color: #c7c7c7; }

.term-fg1.term-fg30 { color: #676767; }
.term-fg1.term-fg31 { color: #ff6d67; }
.term-fg1.term-fg32 { color: #5ff967; }
.term-fg1.term-fg33 { color: #fefb67; }
.term-fg1.term-fg34 { color: #6871ff; }
.term-fg1.term-fg35 { color: #ff76ff; }
.term-fg1.term-fg36 { color: #5ffdff; }
.term-fg1.term-fg37 { color: #feffff; }
		</style>
	</head>
	<body>
		<pre><code>` + renderStr + `</pre></code>
		` + htmlServerSentEvents + `
	</body>
</html>`
}

////////////////////////////////////////////////////////////////////////////////

type ServeOptions struct {
	WarmUpFlag bool
	Dev        chan DevError
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
		bstr, err := os.ReadFile(filepath.Join(RETRO_OUT_DIR, "index.html"))
		if err != nil {
			return err
		}
		contents = strings.Replace(string(bstr), "</body>", fmt.Sprintf("\t%s\n\t</body>", htmlServerSentEvents), 1)
	}

	var (
		msg    = <-options.Dev
		logMsg string
	)

	// Path for HTML and non-HTML resources
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log to stdout
		var nextLogMsg string
		if msg.IsDirty() {
			nextLogMsg = msg.String()
		} else {
			// nextLogMsg = buildServeSuccessString(a.getPort(), msg.dur) // FIXME
			nextLogMsg = buildServeSuccessString(a.getPort(), 0) // TODO
		}
		if logMsg != nextLogMsg {
			logMsg = nextLogMsg
			terminal.Clear(os.Stdout)
			fmt.Println(logMsg)
		}
		// Log to the browser and eagerly return
		if msg.IsDirty() {
			fmt.Fprintln(w, msg.HTML())
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
		http.ServeFile(w, r, filepath.Join(RETRO_OUT_DIR, "index.html"))
	})

	// Paths for dev events and errors
	//
	// - /__dev__
	// - /__err__
	//
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
				case msg = <-options.Dev:
					fmt.Fprint(w, "event: reload\ndata\n\n")
					flusher.Flush()
				case <-r.Context().Done():
					return
				}
			}
		})
		http.HandleFunc("/__err__", func(w http.ResponseWriter, r *http.Request) {
			bstr, err := io.ReadAll(r.Body)
			must(err)
			options.Dev <- newRuntimeError(string(bstr))
		})
	}

	// Log to stdout
	var nextLogMsg string
	if msg.IsDirty() {
		nextLogMsg = msg.String()
	} else {
		// nextLogMsg = buildServeSuccessString(a.getPort(), msg.dur) // FIXME
		nextLogMsg = buildServeSuccessString(a.getPort(), 0) // TODO
	}
	if logMsg != nextLogMsg {
		logMsg = nextLogMsg
		terminal.Clear(os.Stdout)
		fmt.Println(logMsg)
	}

	port := a.getPort()
	for {
		// FIXME: Go doesn't error on used ports?
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
