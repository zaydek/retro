package format

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zaydek/retro/go/pkg/terminal"
)

var accentRegex = regexp.MustCompile("`([^`]+)`")

func Stdout(x interface{}) string {
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
	for arrIndex, el := range arr {
		if el == "" {
			continue
		}
		arr[arrIndex] =
			accentRegex.ReplaceAllString(
				strings.ReplaceAll(el, "\t", "  "), // Tabs -> spaces
				terminal.Cyan("`$1`"),              // Accent
			)
	}
	return strings.Join(arr, "\n")
}

func Stderr(x interface{}) string {
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
	for arrIndex, el := range arr {
		if el == "" {
			continue
		}
		arr[arrIndex] =
			terminal.BoldRed("error:") + " " +
				accentRegex.ReplaceAllString(
					strings.ReplaceAll(el, "\t", "  "), // Tabs -> spaces
					terminal.Magenta("`$1`"),           // Accent
				)
	}
	return strings.Join(arr, "\n")
}

func StderrIPC(str string) string {
	var ret string
	split := strings.Split(strings.TrimRight(str, "\n"), "\n")
	for lineIndex, line := range split {
		if lineIndex > 0 {
			ret += "\n"
		}
		ret +=
			fmt.Sprintf("%s  %s",
				terminal.BoldRed("stderr"),
				line,
			)
	}
	return ret
}
