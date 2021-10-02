package retro

import (
	"fmt"
	"strings"

	"github.com/zaydek/retro/go/pkg/terminal"
)

// (Node.js) stdout  ...
func formatStdoutLine(line string) string {
	stdout := fmt.Sprintf(
		"%s %s  %s",
		terminal.Dim("(Node.js)"),
		terminal.BoldCyan("stdout"),
		line,
	)
	return stdout
}

// (Node.js) stderr  ...
func formatStderrLine(line string) string {
	stdout := fmt.Sprintf(
		"%s %s  %s",
		terminal.Dim("(Node.js)"),
		terminal.BoldRed("stderr"),
		line,
	)
	return stdout
}

// (Node.js) stdout  ...
// (Node.js) stdout  ...
func formatStdoutText(text string) string {
	var stdout string
	split := strings.Split(strings.TrimRight(text, "\n"), "\n")
	for lineIndex, line := range split {
		if lineIndex > 0 {
			stdout += "\n"
		}
		stdout += formatStdoutLine(line)
	}
	return stdout
}

// (Node.js) stderr  ...
// (Node.js) stderr  ...
func formatStderrText(text string) string {
	var stderr string
	split := strings.Split(strings.TrimRight(text, "\n"), "\n")
	for lineIndex, line := range split {
		if lineIndex > 0 {
			stderr += "\n"
		}
		stderr += formatStderrLine(line)
	}
	return stderr
}
