package stdio_logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zaydek/retro/pkg/terminal"
)

type LoggerOptions struct {
	Datetime bool
	Date     bool
	Time     bool
}

type StdioLogger struct {
	format string
}

func extractFormat(opt LoggerOptions) string {
	var format string
	if opt.Datetime {
		format += "Jan 02 03:04:05.000 PM"
	} else {
		if opt.Date {
			format += "Jan 02"
		}
		if opt.Time {
			if format != "" {
				format += " "
			}
			format += "03:04:05.000 PM"
		}
	}
	return format
}

func New(args ...LoggerOptions) *StdioLogger {
	opt := LoggerOptions{}
	if len(args) == 1 {
		opt = args[0]
	}
	logger := &StdioLogger{format: extractFormat(opt)}
	return logger
}

func (l *StdioLogger) Set(opt LoggerOptions) {
	l.format = extractFormat(opt)
}

func (l *StdioLogger) Stdout(args ...interface{}) {
	var tstr string
	if l.format != "" {
		tstr += terminal.Dim(time.Now().Format(l.format))
		tstr += "  "
	}

	str := strings.TrimRight(fmt.Sprint(args...), "\n")
	lines := strings.Split(str, "\n")
	for x, line := range lines {
		lines[x] = fmt.Sprintf("%s%s\x1b[0m", tstr, line)
	}
	fmt.Fprintln(os.Stdout, strings.Join(lines, "\n"))
}

func (l *StdioLogger) Stderr(args ...interface{}) {
	var tstr string
	if l.format != "" {
		tstr += terminal.Dim(time.Now().Format(l.format))
		tstr += "  "
	}

	str := strings.TrimRight(fmt.Sprint(args...), "\n")
	lines := strings.Split(str, "\n")
	for x, line := range lines {
		lines[x] = fmt.Sprintf("%s%s %s\x1b[0m", tstr, terminal.BoldRed("stderr"), line)
	}
	fmt.Fprintln(os.Stderr, strings.Join(lines, "\n"))
}

////////////////////////////////////////////////////////////////////////////////

var stdio = New(LoggerOptions{Datetime: true})

func Set(opt LoggerOptions) {
	stdio.Set(opt)
}

func Stdout(args ...interface{}) {
	stdio.Stdout(args...)
}

func Stderr(args ...interface{}) {
	stdio.Stderr(args...)
}
