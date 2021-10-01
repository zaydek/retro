package ipc

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// Starts a long-lived IPC process. stdout messages are read line-by-line
// whereas stderr messages are read once.
func NewCommand(commandArgs ...string) (stdin, stdout, stderr chan string, err error) {
	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)

	// Get pipes
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		returnError := fmt.Errorf("cmd.StdinPipe: %w", err)
		return nil, nil, nil, returnError
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		returnError := fmt.Errorf("cmd.StdoutPipe: %w", err)
		return nil, nil, nil, returnError
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		returnError := fmt.Errorf("cmd.StderrPipe: %w", err)
		return nil, nil, nil, returnError
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		returnError := fmt.Errorf("cmd.Start: %w", err)
		return nil, nil, nil, returnError
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
	}()

	stderr = make(chan string)
	go func() {
		defer func() {
			stderrPipe.Close()
			close(stderr)
		}()
		// Scan once
		scanner := bufio.NewScanner(stderrPipe)
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			return len(data), data, nil
		})
		scanner.Scan()
		if text := scanner.Text(); text != "" {
			stderr <- strings.TrimRight(
				text,
				"\n", // Remove the EOF
			)
		}
	}()

	return stdin, stdout, stderr, nil
}
