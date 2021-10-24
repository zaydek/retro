package format

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zaydek/retro/go/pkg/terminal"
)

var quoteRegex = regexp.MustCompile("'([^']+)'")

func Stdout(x interface{}) string {
	var out string
	switch v := x.(type) {
	case string:
		out = v
	case error:
		out = v.Error()
	default:
		panic("Internal error")
	}
	arr := strings.Split(out, "\n")
	for arrIndex, str := range arr {
		if str != "" {
			arr[arrIndex] = quoteRegex.ReplaceAllString(
				strings.ReplaceAll(str, "\t", "  "), // Tabs -> spaces
				terminal.Cyan("'$1'"),               // Accent
			)
		}
	}
	return strings.Join(arr, "\n")
}

func Stderr(x interface{}) string {
	var out string
	switch v := x.(type) {
	case string:
		out = v
	case error:
		out = v.Error()
	default:
		panic("Internal error")
	}
	arr := strings.Split(out, "\n")
	for arrIndex, str := range arr {
		if str != "" {
			arr[arrIndex] = quoteRegex.ReplaceAllString(
				strings.ReplaceAll(str, "\t", "  "), // Tabs -> spaces
				terminal.Magenta("'$1'"),            // Accent
			)
			if arrIndex == 0 {
				arr[0] = terminal.Bold(terminal.Red("error:") + " " + arr[0])
			}
		}
	}
	return strings.Join(arr, "\n")
}

func StderrIPC(str string) string {
	var out string
	split := strings.Split(strings.TrimRight(str, "\n"), "\n")
	for lineIndex, line := range split {
		if lineIndex > 0 {
			out += "\n"
		}
		out += fmt.Sprintf("%s  %s",
			terminal.BoldRed("stderr"),
			line,
		)
	}
	return out
}
