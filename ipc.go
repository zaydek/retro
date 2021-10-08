package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
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
		defer func() {
			stdinPipe.Close()
			close(stdin)
		}()
		for arg := range stdin {
			fmt.Fprintln(stdinPipe, arg)
		}
	}()

	go func() {
		defer func() {
			stdoutPipe.Close()
			close(stdout)
		}()
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
		defer func() {
			stderrPipe.Close()
			close(stderr)
		}()
		scanner := bufio.NewScanner(stderrPipe)
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			return len(data), data, nil
		})
		scanner.Scan()
		// NOTE: As of Sass v1.42.1, deprecation warnings are uncontrollable. To
		// suppress these warnings, add a micro-delay and guard for deprecation
		// warnings.
		time.Sleep(50 * time.Millisecond)
		if text := scanner.Text(); text != "" {
			if !strings.HasPrefix(text, "DEPRECATION WARNING") {
				stderr <- strings.TrimRight(text, "\n")
			}
		}
		must(scanner.Err())
	}()

	if err := cmd.Start(); err != nil {
		return nil, nil, nil, err
	}
	return stdin, stdout, stderr, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	stdin, stdout, stderr, err := NewPersistentCommand(ctx, "node", "esbuild.js")
	if err != nil {
		cancel()
		must(err)
	}
	defer cancel()

	stdin <- "build"
	select {
	case outStr := <-stdout:
		fmt.Printf("stdout %s\n", outStr)
	case errStr := <-stderr:
		fmt.Printf("stderr %s\n", errStr)
	}

	stdin <- "build"
	select {
	case outStr := <-stdout:
		fmt.Printf("stdout %s\n", outStr)
	case errStr := <-stderr:
		fmt.Printf("stderr %s\n", errStr)
	}
}
