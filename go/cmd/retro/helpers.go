package retro

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/zaydek/retro/go/pkg/sys"
	"github.com/zaydek/retro/go/pkg/terminal"
)

func backtick(str string) string {
	return "`" + str + "`"
}

////////////////////////////////////////////////////////////////////////////////

// Gets the filesystem path for a URL (adds `index.html` and `.html`) Note that
// `getFilesystemPath` is inverse to `getCanonicalBrowserPath`.
func getFilesystemPath(url string) string {
	ret := url
	if strings.HasSuffix(url, "/") {
		ret += "index.html"
	} else if strings.HasSuffix(url, "/index") {
		ret += ".html"
	} else if ext := filepath.Ext(url); ext == "" {
		ret += ".html"
	}
	return ret
}

// Gets the canonical browser path for a URL (removes `index.html` and `.html`).
// Note that `getCanonicalBrowserPath` is inverse to `getFilesystemPath`.
func getCanonicalBrowserPath(url string) string {
	ret := url
	if strings.HasSuffix(url, "/index.html") {
		ret = ret[:len(ret)-len("index.html")]
	} else if strings.HasSuffix(url, "/index") {
		ret = ret[:len(ret)-len("index")]
	} else if strings.HasSuffix(url, ".html") {
		ret = ret[:len(ret)-len(".html")]
	}
	return ret
}

////////////////////////////////////////////////////////////////////////////////

// Like `__dirname` in Node.js; gets the executable's directory
func getDirname() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	// Edge-case for local development: When running `go run main-*.go`, get the
	// working directory. `main-*` works as heuristic because the entry point
	// filenames are `main_create_retro_app.go` and `main_retro.go`.
	if strings.HasPrefix(filepath.Base(executable), "main_") {
		return os.Getwd()
	}
	// Follows symlinks; get `node_modules/.bin/@zaydek/bin/retro` not
	// `node_modules/.bin/retro`
	return filepath.EvalSymlinks(executable)
}

////////////////////////////////////////////////////////////////////////////////

// https://stackoverflow.com/a/37382208
func getIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

func makeServeSuccess(port int) string {
	ip, err := getIP()
	isOffline := err != nil && strings.HasSuffix(
		err.Error(),
		"dial udp 8.8.8.8:80: connect: network is unreachable",
	)

	wd, _ := os.Getwd()
	base := filepath.Base(wd)
	if isOffline {
		return terminal.Green("Compiled successfully!") + `

You can now view ` + terminal.Bold(base) + ` in the browser.

  ` + terminal.Bold("Local:") + `            ` + fmt.Sprintf("http://localhost:%s", terminal.Bold(port)) + `

Note that the development build is not optimized.
To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.`
	}

	return terminal.Green("Compiled successfully!") + `

You can now view ` + terminal.Bold(base) + ` in the browser.

  ` + terminal.Bold("Local:") + `            ` + fmt.Sprintf("http://localhost:%s", terminal.Bold(port)) + `
  ` + terminal.Bold("On Your Network:") + `  ` + fmt.Sprintf("http://%s:%s", ip, terminal.Bold(port)) + `

Note that the development build is not optimized.
To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.`
}

////////////////////////////////////////////////////////////////////////////////

var epoch = time.Now()

func makeBuildSuccess(directory string) (string, error) {
	var str string

	lsInfos, err := sys.List(directory)
	if err != nil {
		return "", err
	}
	sort.Sort(lsInfos)

	var sum int64
	for _, lsInfo := range lsInfos {
		var color = terminal.Normal
		if strings.HasSuffix(lsInfo.Path, ".html") {
			color = terminal.Normal
		} else if strings.HasSuffix(lsInfo.Path, ".js") || strings.HasSuffix(lsInfo.Path, ".js.map") {
			color = terminal.Yellow
		} else if strings.HasSuffix(lsInfo.Path, ".css") || strings.HasSuffix(lsInfo.Path, ".css.map") {
			color = terminal.Cyan
		} else {
			color = terminal.Dim
		}
		str += fmt.Sprintf("%v%s%v\n", color(lsInfo.Path), strings.Repeat(" ", 40-len(lsInfo.Path)),
			terminal.Dimf("(%s)", sys.ByteCountIEC(lsInfo.Size)))
		if !strings.HasSuffix(lsInfo.Path, ".map") {
			sum += lsInfo.Size
		}
	}

	str += fmt.Sprintln(strings.Repeat(" ", 40) + terminal.Dimf("(%s sum)", sys.ByteCountIEC(sum)))
	str += fmt.Sprintln()
	str += fmt.Sprintln(terminal.Dimf("(%s)", time.Since(epoch)))

	return str, nil
}
