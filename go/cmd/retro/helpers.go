package retro

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/zaydek/retro/go/pkg/terminal"
)

func backtick(str string) string {
	return "`" + str + "`"
}

////////////////////////////////////////////////////////////////////////////////

// Gets the filesystem path for a URL (adds `index.html` and `.html`) Note that
// `getFilesystemPath` is inverse to `getCanonicalBrowserPath`.
func getFilesystemPath(url string) string {
	returnUrl := url
	if strings.HasSuffix(url, "/") {
		returnUrl += "index.html"
	} else if strings.HasSuffix(url, "/index") {
		returnUrl += ".html"
	} else if extension := filepath.Ext(url); extension == "" {
		returnUrl += ".html"
	}
	return returnUrl
}

// // Gets the canonical browser path for a URL (removes `index.html` and `.html`).
// // Note that `getCanonicalBrowserPath` is inverse to `getFilesystemPath`.
// func getCanonicalBrowserPath(url string) string {
// 	returnUrl := url
// 	if strings.HasSuffix(url, "/index.html") {
// 		returnUrl = returnUrl[:len(returnUrl)-len("index.html")]
// 	} else if strings.HasSuffix(url, "/index") {
// 		returnUrl = returnUrl[:len(returnUrl)-len("index")]
// 	} else if strings.HasSuffix(url, ".html") {
// 		returnUrl = returnUrl[:len(returnUrl)-len(".html")]
// 	}
// 	return returnUrl
// }

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
		return nil, fmt.Errorf("net.Dial: %w", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// TODO: Add error-handling
func buildSuccess(port int) error {
	// Get the working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}

	// Get the IP address and a flag describing the user's offline
	ip, err := getIP()
	isNetworkUnreachable := err != nil && strings.HasSuffix(
		err.Error(),
		"dial udp 8.8.8.8:80: connect: network is unreachable",
	)
	if err != nil && !isNetworkUnreachable {
		return fmt.Errorf("getIP: %w", err)
	}

	// Log success message; depends on network access
	base := filepath.Base(cwd)
	if isNetworkUnreachable {
		terminal.Clear(os.Stdout)
		fmt.Println(terminal.Green("Compiled successfully!") + `

You can now view ` + terminal.Bold(base) + ` in the browser.

  ` + terminal.Bold("Local:") + `            ` + fmt.Sprintf("http://localhost:%s", terminal.Bold(port)) + `

Note that the development build is not optimized.
To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.` + "\n")
	} else {
		terminal.Clear(os.Stdout)
		fmt.Println(terminal.Green("Compiled successfully!") + `

You can now view ` + terminal.Bold(base) + ` in the browser.

  ` + terminal.Bold("Local:") + `            ` + fmt.Sprintf("http://localhost:%s", terminal.Bold(port)) + `
  ` + terminal.Bold("On Your Network:") + `  ` + fmt.Sprintf("http://%s:%s", ip, terminal.Bold(port)) + `

Note that the development build is not optimized.
To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.` + "\n")
	}

	return nil
}
