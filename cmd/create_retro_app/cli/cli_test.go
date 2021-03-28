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

func TestParseCommand(t *testing.T) {
	var (
		cmd Command
		err error
	)

	cmd, err = ParseCommand(".")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: ".",
	})

	cmd, err = ParseCommand(".", "--template=js")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: ".",
	})

	cmd, err = ParseCommand(".", "--template=jsx")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: ".",
	})

	cmd, err = ParseCommand(".", "--template=javascript")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: ".",
	})

	cmd, err = ParseCommand(".", "--template=ts")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: ".",
	})

	cmd, err = ParseCommand(".", "--template=tsx")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: ".",
	})

	cmd, err = ParseCommand(".", "--template=typescript")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: ".",
	})

	//////////////////////////////////////////////////////////////////////////////

	cmd, err = ParseCommand("app")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: "app",
	})

	cmd, err = ParseCommand("app", "--template=js")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: "app",
	})

	cmd, err = ParseCommand("app", "--template=jsx")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: "app",
	})

	cmd, err = ParseCommand("app", "--template=javascript")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: "app",
	})

	cmd, err = ParseCommand("app", "--template=ts")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: "app",
	})

	cmd, err = ParseCommand("app", "--template=tsx")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: "app",
	})

	cmd, err = ParseCommand("app", "--template=typescript")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: "app",
	})
}
