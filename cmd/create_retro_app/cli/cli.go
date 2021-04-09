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
		return fmt.Sprintf("Unsupported argument '%s'.", e.BadArgument)
	case BadTemplateValue:
		return "'--template' must be a 'starter', 'sass', or 'mdx' (default 'starter')."
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
		Template:  "starter",
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
			case "starter":
				cmd.Template = "starter"
			case "sass":
				cmd.Template = "sass"
			case "mdx":
				cmd.Template = "mdx"
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
		return Command{}, ErrUsage
	}

	var (
		cmd Command
		err error
	)

	// TODO: Previously --port was not passed as an option to the dev server. It’s
	// not clear whether this is because of os.Args[2:] or something else.
	if cmdArg := os.Args[1]; cmdArg == "version" || cmdArg == "--version" || cmdArg == "-v" {
		return Command{}, ErrVersion
	} else if cmdArg == "usage" || cmdArg == "--usage" {
		return Command{}, ErrUsage
	} else if cmdArg == "help" || cmdArg == "--help" {
		return Command{}, ErrHelp
	} else {
		cmd, err = ParseCommand(os.Args[1:]...)
	}
	return cmd, err
}
