package cli

import (
	"testing"

	"github.com/zaydek/retro/go/pkg/expect"
)

func check(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Fatalf("check: %s", err)
}

func TestDevCommand(t *testing.T) {
	var (
		command DevCommand
		err     error
	)

	command, err = ParseDevCommand()
	check(t, err)
	expect.DeepEqual(t, command, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	command, err = ParseDevCommand("--port=8000")
	check(t, err)
	expect.DeepEqual(t, command, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	command, err = ParseDevCommand("--port=3000")
	check(t, err)
	expect.DeepEqual(t, command, DevCommand{
		Port:      3000,
		Sourcemap: true,
	})

	command, err = ParseDevCommand("--sourcemap")
	check(t, err)
	expect.DeepEqual(t, command, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	command, err = ParseDevCommand("--sourcemap=true")
	check(t, err)
	expect.DeepEqual(t, command, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	command, err = ParseDevCommand("--sourcemap=false")
	check(t, err)
	expect.DeepEqual(t, command, DevCommand{
		Port:      8000,
		Sourcemap: false,
	})
}

func TestBuildCommand(t *testing.T) {
	var (
		command BuildCommand
		err     error
	)

	command, err = ParseBuildCommand()
	check(t, err)
	expect.DeepEqual(t, command, BuildCommand{
		Sourcemap: true,
	})

	command, err = ParseBuildCommand("--sourcemap")
	check(t, err)
	expect.DeepEqual(t, command, BuildCommand{
		Sourcemap: true,
	})

	command, err = ParseBuildCommand("--sourcemap=true")
	check(t, err)
	expect.DeepEqual(t, command, BuildCommand{
		Sourcemap: true,
	})

	command, err = ParseBuildCommand("--sourcemap=false")
	check(t, err)
	expect.DeepEqual(t, command, BuildCommand{
		Sourcemap: false,
	})
}

func TestServeCommand(t *testing.T) {
	var (
		command ServeCommand
		err     error
	)

	command, err = ParseServeCommand()
	check(t, err)
	expect.DeepEqual(t, command, ServeCommand{
		Port: 8000,
	})

	command, err = ParseServeCommand("--port=8000")
	check(t, err)
	expect.DeepEqual(t, command, ServeCommand{
		Port: 8000,
	})

	command, err = ParseServeCommand("--port=3000")
	check(t, err)
	expect.DeepEqual(t, command, ServeCommand{
		Port: 3000,
	})
}
