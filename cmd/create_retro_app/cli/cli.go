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
	VersionError = errors.New("cli: version error")
	UsageError   = errors.New("cli: usage error")
	HelpError    = errors.New("cli: help error")
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
		return fmt.Sprintf("Unsupported argument '%s'.", e.BadArgument)
	case BadTemplateValue:
		return fmt.Sprintf("'--template' must be a 'javascript' or 'typescript' (default 'javascript').")
	case BadDirectoryValue:
		cwd, _ := os.Getwd()
		return fmt.Sprintf("Use '.' explicitly to use '%s'.", filepath.Join("..", filepath.Base(cwd)))
	}
	panic("Internal error")
}

func (e CommandError) Unwrap() error {
	return e.Err
}

func ParseCommand(args ...string) (Command, error) {
	cmd := Command{
		Template:  "javascript",
		Directory: "",
	}
	var once sync.Once
	for _, arg := range args {
		err := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--template") {
			if len(arg) <= len("--template=") {
				return Command{}, CommandError{Kind: BadTemplateValue}
			}
			switch strings.ToLower(arg[len("--template="):]) {
			case "js":
				fallthrough
			case "jsx":
				fallthrough
			case "javascript":
				cmd.Template = "javascript"
			case "ts":
				fallthrough
			case "tsx":
				fallthrough
			case "typescript":
				cmd.Template = "typescript"
			default:
				return Command{}, CommandError{Kind: BadTemplateValue}
			}
		} else if !strings.HasPrefix(arg, "--") {
			once.Do(func() {
				cmd.Directory = arg
			})
		} else {
			return Command{}, err
		}
	}
	if cmd.Directory == "" {
		return Command{}, CommandError{Kind: BadDirectoryValue}
	}
	return cmd, nil
}

func ParseCLIArguments() (Command, error) {
	if len(os.Args) < 2 {
		return Command{}, UsageError
	}

	var (
		cmd Command
		err error
	)

	if cmdArg := os.Args[1]; cmdArg == "version" || cmdArg == "--version" || cmdArg == "-v" {
		return Command{}, VersionError
	} else if cmdArg == "usage" || cmdArg == "--usage" {
		return Command{}, UsageError
	} else if cmdArg == "help" || cmdArg == "--help" {
		return Command{}, HelpError
	} else {
		cmd, err = ParseCommand(os.Args[1:]...)
	}
	return cmd, err
}
