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
	BadTemplateValue
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
	case BadTemplateValue:
		return "`--template` must be a `starter`, `sass`, or `mdx` (default `starter`)."
	case BadDirectoryValue:
		cwd, _ := os.Getwd()
		return fmt.Sprintf("Use `.` explicitly to use `%s`.", filepath.Join("..", filepath.Base(cwd)))
	}
	panic("Internal error")
}

func (e CommandError) Unwrap() error {
	return e.Err
}

func ParseCommand(args ...string) (*CreateCommand, error) {
	command := &CreateCommand{
		Template:  "starter",
		Directory: "",
	}
	var once sync.Once
	for _, arg := range args {
		err := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--template") {
			if len(arg) <= len("--template=") {
				return nil, CommandError{Kind: BadTemplateValue}
			}
			switch strings.ToLower(arg[len("--template="):]) {
			case "starter":
				command.Template = "starter"
			case "sass":
				command.Template = "sass"
			case "mdx":
				command.Template = "mdx"
			default:
				return nil, CommandError{Kind: BadTemplateValue}
			}
		} else if !strings.HasPrefix(arg, "--") {
			once.Do(func() {
				command.Directory = arg
			})
		} else {
			return nil, err
		}
	}
	if command.Directory == "" {
		return nil, CommandError{Kind: BadDirectoryValue}
	}
	return command, nil
}

func ParseCLIArguments() (*CreateCommand, error) {
	if len(os.Args) < 2 {
		return nil, ErrUsage
	}

	var (
		command *CreateCommand
		err     error
	)

	// TODO: Previously --port was not passed as an option to the dev server. It’s
	// not clear whether this is because of os.Args[2:] or something else.
	if cmdArg := os.Args[1]; cmdArg == "version" || cmdArg == "--version" || cmdArg == "-v" {
		return nil, ErrVersion
	} else if cmdArg == "usage" || cmdArg == "--usage" {
		return nil, ErrUsage
	} else if cmdArg == "help" || cmdArg == "--help" {
		return nil, ErrHelp
	} else {
		command, err = ParseCommand(os.Args[1:]...)
	}
	return command, err
}
