package cli

import (
	"testing"

	"github.com/zaydek/retro/pkg/expect"
)

func must(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Fatal(err)
}

func TestDevCommand(t *testing.T) {
	var cmd DevCommand
	var err error

	cmd, err = ParseDevCommand()
	must(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	cmd, err = ParseDevCommand("--port=8000")
	must(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	cmd, err = ParseDevCommand("--port=3000")
	must(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      3000,
		Sourcemap: true,
	})

	cmd, err = ParseDevCommand("--sourcemap")
	must(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	cmd, err = ParseDevCommand("--sourcemap=true")
	must(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      8000,
		Sourcemap: true,
	})

	cmd, err = ParseDevCommand("--sourcemap=false")
	must(t, err)
	expect.DeepEqual(t, cmd, DevCommand{
		Port:      8000,
		Sourcemap: false,
	})
}

func TestBuildCommand(t *testing.T) {
	var cmd BuildCommand
	var err error

	cmd, err = ParseExportCommand()
	must(t, err)
	expect.DeepEqual(t, cmd, BuildCommand{
		Sourcemap: true,
	})

	cmd, err = ParseExportCommand("--sourcemap")
	must(t, err)
	expect.DeepEqual(t, cmd, BuildCommand{
		Sourcemap: true,
	})

	cmd, err = ParseExportCommand("--sourcemap=true")
	must(t, err)
	expect.DeepEqual(t, cmd, BuildCommand{
		Sourcemap: true,
	})

	cmd, err = ParseExportCommand("--sourcemap=false")
	must(t, err)
	expect.DeepEqual(t, cmd, BuildCommand{
		Sourcemap: false,
	})
}

func TestServeCommand(t *testing.T) {
	var cmd ServeCommand
	var err error

	cmd, err = ParseServeCommand()
	must(t, err)
	expect.DeepEqual(t, cmd, ServeCommand{
		Port: 8000,
	})

	cmd, err = ParseServeCommand("--port=8000")
	must(t, err)
	expect.DeepEqual(t, cmd, ServeCommand{
		Port: 8000,
	})

	cmd, err = ParseServeCommand("--port=3000")
	must(t, err)
	expect.DeepEqual(t, cmd, ServeCommand{
		Port: 3000,
	})
}
