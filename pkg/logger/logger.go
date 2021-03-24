package logger

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/zaydek/retro/pkg/terminal"
)

var accentRegex = regexp.MustCompile(`('[^']+')`)

func Transform(str string, accent func(...interface{}) string) string {
	arr := strings.Split(strings.TrimRight(str, "\n"), "\n")
	for x := range arr {
		if arr[x] == "" {
			continue
		}
		arr[x] = strings.ReplaceAll(arr[x], "\t", "  ")
		arr[x] = accentRegex.ReplaceAllString(arr[x], accent("$1"))
		arr[x] = " " + arr[x]
	}
	return strings.Join(arr, "\n")
}

func OK(str string) {
	out := terminal.Boldf(" > %s%s\n", terminal.Green("ok:"), Transform(str, terminal.Green))
	fmt.Fprintln(os.Stdout, out)
}

func Warning(err error) {
	out := terminal.Boldf(" > %s%s\n", terminal.Yellow("warning:"), Transform(err.Error(), terminal.Magenta))
	fmt.Fprintln(os.Stderr, out)
}

// func SourceError(source string, line, column int, err error) {
// 	out := terminal.Boldf(" > %s: %s%s\n", fmt.Sprintf("%s:%d:%d", source, line, column), terminal.Red("error:"), Transform(err.Error(), terminal.Magenta))
// 	fmt.Fprintln(os.Stderr, out)
// }

func FatalError(err error) {
	out := terminal.Boldf(" > %s%s\n", terminal.Red("error:"), Transform(err.Error(), terminal.Magenta))
	fmt.Fprintln(os.Stderr, out)
	os.Exit(1)
}
