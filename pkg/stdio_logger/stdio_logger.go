package stdio_logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zaydek/retro/pkg/terminal"
)

type LoggerOptions struct {
	Datetime bool // Mar 02 15:04:05 AM
	Date     bool // Mar 02
	Time     bool // 15:04:05 AM
}

type StdioLogger struct {
	fmt string
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
	logger := &StdioLogger{fmt: extractFormat(opt)}
	return logger
}

func (l *StdioLogger) Set(opt LoggerOptions) {
	l.fmt = extractFormat(opt)
}

func (l *StdioLogger) TransformStdout(args ...interface{}) string {
	var tstr string
	if l.fmt != "" {
		tstr += terminal.Dim(time.Now().Format(l.fmt))
	}
	arr := strings.Split(strings.TrimRight(fmt.Sprint(args...), "\n"), "\n")
	for x, v := range arr {
		arr[x] = fmt.Sprintf("%s  %s\x1b[0m", tstr, v)
	}
	return strings.Join(arr, "\n")
}

func (l *StdioLogger) TransformStderr(args ...interface{}) string {
	var tstr string
	if l.fmt != "" {
		tstr += terminal.Dim(time.Now().Format(l.fmt))
		tstr += " " // Add space for 'stderr'
	}
	arr := strings.Split(strings.TrimRight(fmt.Sprint(args...), "\n"), "\n")
	for x, v := range arr {
		arr[x] = fmt.Sprintf("%s%s  %s\x1b[0m", tstr, terminal.BoldRed("stderr"), v)
	}
	return strings.Join(arr, "\n")
}

func (l *StdioLogger) Stdout(args ...interface{}) {
	fmt.Fprintln(os.Stdout, l.TransformStdout(args...))
}

func (l *StdioLogger) Stderr(args ...interface{}) {
	fmt.Fprintln(os.Stderr, l.TransformStderr(args...))
}

////////////////////////////////////////////////////////////////////////////////

var stdio = New(LoggerOptions{Datetime: true})

func Set(opt LoggerOptions) {
	stdio.Set(opt)
}

func TransformStdout(args ...interface{}) string {
	return stdio.TransformStdout(args...)
}

func TransformStderr(args ...interface{}) string {
	return stdio.TransformStderr(args...)
}

func Stdout(args ...interface{}) {
	stdio.Stdout(args...)
}

func Stderr(args ...interface{}) {
	stdio.Stderr(args...)
}
