package v8

import (
	"regexp"
	"strconv"
	"strings"
)

type StackFrame struct {
	Caller string // E.g. Object.<anonymous>
	Source string // E.g. foo/bar/baz.ext
	Line   int    // E.g. 1
	Column int    // E.g. 2
}

type StackTrace struct {
	Error  string
	Frames []StackFrame
}

var frameRe = regexp.MustCompile(`^    at ([^(]+) \((.*):(\d+):(\d+)\)$`)

func parseFrame(frameStr string) StackFrame {
	matches := frameRe.FindAllStringSubmatch(frameStr, -1)
	if matches == nil || len(matches[0]) != 5 {
		panic("Internal error")
	}

	lno, _ := strconv.Atoi(matches[0][3])
	cno, _ := strconv.Atoi(matches[0][4])

	frame := StackFrame{
		Caller: matches[0][1],
		Source: matches[0][2],
		Line:   lno,
		Column: cno,
	}
	return frame
}

func NewStackTrace(stackStr string) StackTrace {
	var trace StackTrace

	if !strings.HasPrefix(stackStr, "Error: ") {
		panic("Internal error")
	}

	// Remove WS from the end of the stack
	arr := strings.Split(strings.TrimSpace(stackStr), "\n")
	for x, v := range arr {
		if x == 0 {
			trace.Error += v[len("Error: "):]
			continue
		}
		if !strings.HasPrefix(v, "    at ") {
			// trace.Error cannot be "" here so we add "\n"
			trace.Error += "\n" + v
		} else {
			trace.Frames = append(trace.Frames, parseFrame(v))
		}
	}
	return trace
}
