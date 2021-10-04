package retro

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zaydek/retro/go/cmd/retro/unix"
	"github.com/zaydek/retro/go/pkg/terminal"
)

func quote(str string) string {
	return "`" + str + "`"
}

////////////////////////////////////////////////////////////////////////////////

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

// func getBrowserPath(url string) string {
// 	ret := url
// 	if strings.HasSuffix(url, "/index.html") {
// 		ret = ret[:len(ret)-len("index.html")]
// 	} else if strings.HasSuffix(url, "/index") {
// 		ret = ret[:len(ret)-len("index")]
// 	} else if strings.HasSuffix(url, ".html") {
// 		ret = ret[:len(ret)-len(".html")]
// 	}
// 	return ret
// }

////////////////////////////////////////////////////////////////////////////////

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

var epoch = time.Now()

func buildBuildSuccessString(dir string) (string, error) {
	var ret string
	ls, err := unix.List(dir)
	if err != nil {
		return "", err
	}
	for _, info := range ls {
		var (
			color = terminal.Dim
			ext   = filepath.Ext(info.Path)
		)
		switch ext {
		case ".html":
			color = terminal.Normal
		case ".js":
			fallthrough
		case ".js.map":
			color = terminal.Yellow
		case ".css":
			fallthrough
		case ".css.map":
			color = terminal.Cyan
		}
		ret += fmt.Sprintf("%v%s%v\n", color(info.Path), strings.Repeat(" ", 40-len(info.Path)),
			terminal.Dim(unix.HumanReadable(info.Size)))
	}
	ret += fmt.Sprintln()
	ret += fmt.Sprintln(terminal.Dimf("%dms", time.Since(epoch).Milliseconds()))
	return ret, nil
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

func buildServeSucessString(port int, d time.Duration) string {
	ip, err := getIP()
	isOffline := err != nil && strings.HasSuffix(err.Error(), "dial udp 8.8.8.8:80: connect: network is unreachable")

	wd, _ := os.Getwd()
	base := filepath.Base(wd)
	if isOffline {
		return terminal.Green("Compiled successfully!") + `

You can now view ` + terminal.Bold(base) + ` in the browser.

  ` + terminal.Bold("Local:") + `            ` + fmt.Sprintf("http://localhost:%s", terminal.Bold(port)) + `

Note that the development build is not optimized.
To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.

` + terminal.Dimf("%dms", d.Milliseconds())
	}

	return terminal.Green("Compiled successfully!") + `

You can now view ` + terminal.Bold(base) + ` in the browser.

  ` + terminal.Bold("Local:") + `            ` + fmt.Sprintf("http://localhost:%s", terminal.Bold(port)) + `
  ` + terminal.Bold("On Your Network:") + `  ` + fmt.Sprintf("http://%s:%s", ip, terminal.Bold(port)) + `

Note that the development build is not optimized.
To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.

` + terminal.Dimf("%dms", d.Milliseconds())
}
