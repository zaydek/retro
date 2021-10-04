package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ErrVersion = errors.New("cli: version error")
	ErrUsage   = errors.New("cli: usage error")
	ErrHelp    = errors.New("cli: help error")
)

type ErrorKind int

const (
	BadArgument ErrorKind = iota
	BadDirectoryValue
)

type CommandError struct {
	Kind        ErrorKind
	BadArgument string
	Err         error
}

func (e CommandError) Error() string {
	switch e.Kind {
	case BadArgument:
		return fmt.Sprintf("Unsupported argument `%s`.", e.BadArgument)
	case BadDirectoryValue:
		wd, _ := os.Getwd()
		return fmt.Sprintf("Use `.` explicitly to use the working directory `%s`.", filepath.Base(wd))
	}
	panic("Internal error")
}

func (e CommandError) Unwrap() error {
	return e.Err
}

func ParseCommand(args ...string) (CreateCommand, error) {
	var command = CreateCommand{}
	var once sync.Once
	for _, arg := range args {
		err := CommandError{Kind: BadArgument, BadArgument: arg}
		if !strings.HasPrefix(arg, "--") {
			once.Do(func() {
				command.Directory = arg
			})
		} else {
			return CreateCommand{}, err
		}
	}
	if command.Directory == "" {
		return CreateCommand{}, CommandError{Kind: BadDirectoryValue}
	}
	return command, nil
}

func ParseCLIArguments() (CreateCommand, error) {
	if len(os.Args) < 2 {
		return CreateCommand{}, ErrUsage
	}

	var (
		command CreateCommand
		err     error
	)

	if cmdArg := os.Args[1]; cmdArg == "version" || cmdArg == "--version" || cmdArg == "-v" {
		return CreateCommand{}, ErrVersion
	} else if cmdArg == "usage" || cmdArg == "--usage" {
		return CreateCommand{}, ErrUsage
	} else if cmdArg == "help" || cmdArg == "--help" {
		return CreateCommand{}, ErrHelp
	} else {
		command, err = ParseCommand(os.Args[1:]...)
	}
	return command, err
}
