package cli

import (
	"testing"

	"github.com/zaydek/retro/pkg/expect"
)

func check(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Fatalf("check: %s", err)
}

func TestDevCommand(t *testing.T) {
	var (
		cmd DevCommand
		err error
	)

	cmd, err = ParseDevCommand()
	check(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	cmd, err = ParseDevCommand("--port=8000")
	check(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	cmd, err = ParseDevCommand("--port=3000")
	check(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      3000,
		Sourcemap: true,
	})

	cmd, err = ParseDevCommand("--sourcemap")
	check(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	cmd, err = ParseDevCommand("--sourcemap=true")
	check(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	cmd, err = ParseDevCommand("--sourcemap=false")
	check(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      8000,
		Sourcemap: false,
	})
}

func TestBuildCommand(t *testing.T) {
	var cmd BuildCommand
	var err error

	cmd, err = ParseBuildCommand()
	check(t, err)
	expect.DeepEqual(t, cmd, BuildCommand{
		Sourcemap: true,
	})

	cmd, err = ParseBuildCommand("--sourcemap")
	check(t, err)
	expect.DeepEqual(t, cmd, BuildCommand{
		Sourcemap: true,
	})

	cmd, err = ParseBuildCommand("--sourcemap=true")
	check(t, err)
	expect.DeepEqual(t, cmd, BuildCommand{
		Sourcemap: true,
	})

	cmd, err = ParseBuildCommand("--sourcemap=false")
	check(t, err)
	expect.DeepEqual(t, cmd, BuildCommand{
		Sourcemap: false,
	})
}

func TestServeCommand(t *testing.T) {
	var cmd ServeCommand
	var err error

	cmd, err = ParseServeCommand()
	check(t, err)
	expect.DeepEqual(t, cmd, ServeCommand{
		Port: 8000,
	})

	cmd, err = ParseServeCommand("--port=8000")
	check(t, err)
	expect.DeepEqual(t, cmd, ServeCommand{
		Port: 8000,
	})

	cmd, err = ParseServeCommand("--port=3000")
	check(t, err)
	expect.DeepEqual(t, cmd, ServeCommand{
		Port: 3000,
	})
}
