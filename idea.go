package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/zaydek/retro/go/pkg/terminal"
)

// var ErrNetworkIsUnreachable = errors.New("dial udp 8.8.8.8:80: connect: network is unreachable")

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
  ` + terminal.Bold("On Your Network:") + `  ` + fmt.Sprintf("http://%s:%s", ip, terminal.Bold(port)) + `

Note that the development build is not optimized.
To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.` + "\n")
	} else {
		terminal.Clear(os.Stdout)
		fmt.Println(terminal.Green("Compiled successfully!") + `

You can now view ` + terminal.Bold(base) + ` in the browser.

  ` + terminal.Bold("Local:") + `            ` + fmt.Sprintf("http://localhost:%s", terminal.Bold(port)) + `

Note that the development build is not optimized.
To create a production build, use ` + terminal.Cyan("npm run build") + ` or ` + terminal.Cyan("yarn build") + `.` + "\n")
	}

	return nil
}

func main() {
	_, err := getIP()
	isNetworkUnreachable := err != nil && strings.HasSuffix(
		err.Error(),
		"dial udp 8.8.8.8:80: connect: network is unreachable",
	)
	if err != nil && !isNetworkUnreachable {
		panic(fmt.Errorf("getIP: %w", err))
	}
	fmt.Printf("isNetworkUnreachable=%v\n", isNetworkUnreachable)
	// fmt.Println(ip)
}
