package cli

import (
	"testing"

	"github.com/zaydek/retro/go/pkg/expect"
)

func must(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Fatal(err)
}

func TestParseCommand(t *testing.T) {
	var (
		command CreateCommand
		err     error
	)

	command, err = ParseCommand(".")
	must(t, err)
	expect.DeepEqual(t, command, CreateCommand{
		Directory: "retro",
	})

	command, err = ParseCommand("app")
	must(t, err)
	expect.DeepEqual(t, command, CreateCommand{
		Directory: "app",
	})
}
