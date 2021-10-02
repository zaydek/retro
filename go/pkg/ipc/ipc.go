package ipc

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Starts a long-lived IPC process. stdout messages are read line-by-line
// whereas stderr messages are read once.
func NewCommand(ctx context.Context, commandArgs ...string) (stdin, stdout, stderr chan string, err error) {
	command := exec.CommandContext(ctx, commandArgs[0], commandArgs[1:]...)

	stdinPipe, err := command.StdinPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	stdoutPipe, err := command.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	stderrPipe, err := command.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	if err := command.Start(); err != nil {
		return nil, nil, nil, err
	}

	stdin = make(chan string)
	go func() {
		defer func() {
			stdinPipe.Close()
			close(stdin)
		}()
		for message := range stdin {
			fmt.Fprintln(stdinPipe, message)
		}
	}()

	stdout = make(chan string)
	go func() {
		defer func() {
			stdoutPipe.Close()
			close(stdout)
		}()
		// Scan line-by-line
		scanner := bufio.NewScanner(stdoutPipe)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
		for scanner.Scan() {
			if line := scanner.Text(); line != "" {
				stdout <- line
			}
		}
		must(scanner.Err())
	}()

	stderr = make(chan string)
	go func() {
		defer func() {
			stderrPipe.Close()
			close(stderr)
		}()
		scanner := bufio.NewScanner(stderrPipe)
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			return len(data), data, nil
		})
		scanner.Scan()
		if text := scanner.Text(); text != "" {
			// Remove the EOF
			stderr <- strings.TrimRight(text, "\n")
		}
		must(scanner.Err())
	}()

	return stdin, stdout, stderr, nil
}
