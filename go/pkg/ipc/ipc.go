package ipc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func NewCommand(args ...string) (string, string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}
	if err := cmd.Start(); err != nil {
		return "", "", err
	}
	stdoutRaw, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return "", "", err
	}
	stderrRaw, err := io.ReadAll(stderrPipe)
	if err != nil {
		return "", "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", "", err
	}
	return string(stdoutRaw), string(stderrRaw), nil
}

func NewPersistentCommand(ctx context.Context, args ...string) (chan string, <-chan string, <-chan string, error) {
	var (
		stdin  = make(chan string)
		stdout = make(chan string)
		stderr = make(chan string)
	)

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	go func() {
		for arg := range stdin {
			fmt.Fprintln(stdinPipe, arg)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
		for scanner.Scan() {
			if line := scanner.Text(); line != "" {
				stdout <- line
			}
		}
		must(scanner.Err())
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			return len(data), data, nil
		})
		scanner.Scan()
		if text := scanner.Text(); text != "" {
			stderr <- strings.TrimRight(text, "\n")
		}
		must(scanner.Err())
	}()

	if err := cmd.Start(); err != nil {
		return nil, nil, nil, err
	}
	return stdin, stdout, stderr, nil
}
