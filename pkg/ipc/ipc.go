package ipc

import (
	"bufio"
	"encoding/json"
	"os/exec"
)

type Request struct {
	Kind string
	Data interface{}
}

type Response struct {
	Kind string
	Data json.RawMessage
}

// NewCommand starts a new IPC command.
func NewCommand(args ...string) (stdin chan Request, stdout chan Response, stderr chan string, err error) {
	cmd := exec.Command(args[0], args[1:]...)

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

	stdin = make(chan Request)
	go func() {
		defer stdinPipe.Close()
		for msg := range stdin {
			bstr, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}
			// Add an EOF so 'await stdin()' can process
			stdinPipe.Write(append(bstr, '\n'))
		}
	}()

	stdout = make(chan Response)
	go func() {
		defer func() {
			stdoutPipe.Close()
			close(stdout)
		}()
		// Increase the buffer
		scanner := bufio.NewScanner(stdoutPipe)
		buf := make([]byte, 1024*1024)
		scanner.Buffer(buf, len(buf))
		for scanner.Scan() {
			var res Response
			if err := json.Unmarshal(scanner.Bytes(), &res); err != nil {
				panic(err)
			}
			stdout <- res
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}()

	stderr = make(chan string)
	go func() {
		defer func() {
			stderrPipe.Close()
			close(stderr)
		}()
		// Read from start-to-end
		// https://golang.org/pkg/bufio/#SplitFunc
		scanner := bufio.NewScanner(stderrPipe)
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			return len(data), data, nil
		})
		for scanner.Scan() {
			// if str := scanner.Text(); str != "" {
			// 	stderr <- str
			// }
			stderr <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}()

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, nil, nil, err
	}

	return stdin, stdout, stderr, nil
}

////////////////////////////////////////////////////////////////////////////////

type Service struct {
	Stdin  chan Request
	Stdout chan Response
	Stderr chan string
}

func (s Service) Send(msg Request, ptr interface{}) (stderr string, err error) {
	s.Stdin <- msg

loop:
	for {
		select {
		case out := <-s.Stdout:
			if out.Kind == "eof" {
				if err := json.Unmarshal(out.Data, ptr); err != nil {
					return "", err
				}
				break loop
			}
		case str := <-s.Stderr:
			stderr = str
		}
	}
	return stderr, nil
}
