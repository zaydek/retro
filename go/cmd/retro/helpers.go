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
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(exe)
	tmpDir := os.TempDir()
	if dir != "" && tmpDir != "" && strings.Contains(dir, tmpDir) {
		wd, _ := os.Getwd()
		return wd, nil
	}
	return dir, nil
}

////////////////////////////////////////////////////////////////////////////////

func buildBuildSuccessString(dir string, dur time.Duration) (string, error) {
	var out string
	ls, err := unix.List(dir)
	if err != nil {
		return "", err
	}
	for _, info := range ls {
		var (
			color = terminal.Dim
			ext   = filepath.Ext(info.Path)
		)
		if info.IsDir() || strings.HasSuffix(ext, ".map") {
			continue
		}
		switch ext {
		case ".html":
			color = terminal.Normal
		case ".css":
			color = terminal.Cyan
		case ".js":
			color = terminal.Yellow
		}
		out += fmt.Sprintf("%v%s%v\n",
			color(info.Path),
			strings.Repeat(" ", 60-len("XXX.X KB")-len(info.Path)),
			terminal.Dim(unix.HumanReadable(info.Size)),
		)
	}
	out += fmt.Sprintln()
	out += fmt.Sprintln(terminal.Dimf("%dms", dur.Milliseconds()))
	return out, nil
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

func buildServeSuccessString(port int, dur time.Duration) string {
	ip, err := getIP()
	isOffline := err != nil && strings.HasSuffix(err.Error(), "dial udp 8.8.8.8:80: connect: network is unreachable")

	wd, _ := os.Getwd()
	base := filepath.Base(wd)
	if isOffline {
		return terminal.Greenf("Compiled successfully! %s", terminal.Dimf("(%s)", os.Getenv("RETRO_V_VERSION"))) + `

You can now view ` + terminal.Bold(base) + ` in the browser.

  ` + terminal.Bold("Local:") + `            ` + fmt.Sprintf("http://localhost:%s", terminal.Bold(port)) + `

Note that the development build is not optimized.
To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.

` + terminal.Dimf("%dms", dur.Milliseconds())
	} else {
		return terminal.Greenf("Compiled successfully! %s", terminal.Dimf("(%s)", os.Getenv("RETRO_V_VERSION"))) + `

You can now view ` + terminal.Bold(base) + ` in the browser.

  ` + terminal.Bold("Local:") + `            ` + fmt.Sprintf("http://localhost:%s", terminal.Bold(port)) + `
  ` + terminal.Bold("On Your Network:") + `  ` + fmt.Sprintf("http://%s:%s", ip, terminal.Bold(port)) + `

Note that the development build is not optimized.
To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.

` + terminal.Dimf("%dms", dur.Milliseconds())
	}
}
