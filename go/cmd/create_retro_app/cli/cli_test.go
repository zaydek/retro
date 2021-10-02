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

func TestParseCommand(t *testing.T) {
	var (
		command *CreateCommand
		err     error
	)

	command, err = ParseCommand(".")
	check(t, err)
	expect.DeepEqual(t, *command, CreateCommand{
		Template:  "starter",
		Directory: ".",
	})

	command, err = ParseCommand(".", "--template=starter")
	check(t, err)
	expect.DeepEqual(t, *command, CreateCommand{
		Template:  "starter",
		Directory: ".",
	})

	command, err = ParseCommand(".", "--template=sass")
	check(t, err)
	expect.DeepEqual(t, *command, CreateCommand{
		Template:  "sass",
		Directory: ".",
	})

	//////////////////////////////////////////////////////////////////////////////

	command, err = ParseCommand("app")
	check(t, err)
	expect.DeepEqual(t, *command, CreateCommand{
		Template:  "starter",
		Directory: "app",
	})

	command, err = ParseCommand("app", "--template=starter")
	check(t, err)
	expect.DeepEqual(t, *command, CreateCommand{
		Template:  "starter",
		Directory: "app",
	})

	command, err = ParseCommand("app", "--template=sass")
	check(t, err)
	expect.DeepEqual(t, *command, CreateCommand{
		Template:  "sass",
		Directory: "app",
	})
}
