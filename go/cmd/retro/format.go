package retro

import (
	"fmt"
	"strings"

	"github.com/zaydek/retro/go/pkg/terminal"
)

// (retro:node) stderr  ...
func formatStderrLine(line string) string {
	stdout := fmt.Sprintf(
		"%s %s  %s",
		terminal.Dim("(retro:node)"),
		terminal.BoldRed("stderr"),
		line,
	)
	return stdout
}

// (retro:node) stderr  ...
// (retro:node) stderr  ...
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
