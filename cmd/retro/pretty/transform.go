package pretty

import (
	"regexp"
	"strings"

	"github.com/zaydek/retro/pkg/terminal"
)

var accentRegex = regexp.MustCompile(`'([^']+)'`)

func Inset(str string) string {
	arr := strings.Split(str, "\n")
	for x, v := range arr {
		arr[x] = " " + v
	}
	return strings.Join(arr, "\n")
}

func Spaces(str string) string {
	arr := strings.Split(str, "\n")
	for x := range arr {
		if arr[x] == "" {
			continue
		}
		arr[x] = strings.ReplaceAll(arr[x], "\t", "  ")
	}
	return strings.Join(arr, "\n")
}

func Accent(str string, accent func(args ...interface{}) string) string {
	arr := strings.Split(str, "\n")
	for x := range arr {
		if arr[x] == "" {
			continue
		}
		arr[x] = accentRegex.ReplaceAllString(arr[x], accent("'$1'"))
	}
	return strings.Join(arr, "\n")
}

func Error(str string) string {
	arr := strings.Split(str, "\n")
	arr[0] = terminal.Boldf("%s %s", terminal.Red("error:"), Accent(arr[0], terminal.Magenta))
	return Spaces(strings.Join(arr, "\n"))
}
