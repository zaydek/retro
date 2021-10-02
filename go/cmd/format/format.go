package format

import (
	"regexp"
	"strings"

	"github.com/zaydek/retro/go/pkg/terminal"
)

func Pad(str string) string {
	arr := strings.Split(str, "\n")
	for x, v := range arr {
		arr[x] = " " + v
	}
	return strings.Join(arr, "\n")
}

func Tabs(str string) string {
	arr := strings.Split(str, "\n")
	for x := range arr {
		if arr[x] == "" {
			continue
		}
		arr[x] = strings.ReplaceAll(arr[x], "\t", "  ")
	}
	return strings.Join(arr, "\n")
}

var accentRe = regexp.MustCompile("`([^`]+)`")

func Accent(str string, accent func(args ...interface{}) string) string {
	arr := strings.Split(str, "\n")
	for x := range arr {
		if arr[x] == "" {
			continue
		}
		arr[x] = accentRe.ReplaceAllString(arr[x], accent("`$1`"))
	}
	return strings.Join(arr, "\n")
}

func NonError(x interface{}) string {
	var str string
	switch v := x.(type) {
	case string:
		str = v
	case error:
		str = v.Error()
	default:
		panic("Internal error")
	}
	return Pad(Tabs(Accent(str, terminal.Cyan)))
}

func Error(x interface{}) string {
	var str string
	switch v := x.(type) {
	case string:
		str = v
	case error:
		str = v.Error()
	default:
		panic("Internal error")
	}
	arr := strings.Split(str, "\n")
	arr[0] = terminal.Boldf("%s %s", terminal.Red("error:"), Accent(arr[0], terminal.Magenta))
	return strings.Join(arr, "\n")
}
