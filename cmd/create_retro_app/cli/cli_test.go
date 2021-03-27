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

func TestParseArguments(t *testing.T) {
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

	//////////////////////////////////////////////////////////////////////////////

	cmd, err = ParseCommand("--template=js")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: ".",
	})

	cmd, err = ParseCommand("--template=jsx")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: ".",
	})

	cmd, err = ParseCommand("--template=javascript")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: ".",
	})

	cmd, err = ParseCommand("--template=ts")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: ".",
	})

	cmd, err = ParseCommand("--template=tsx")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: ".",
	})

	cmd, err = ParseCommand("--template=typescript")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: ".",
	})

	//////////////////////////////////////////////////////////////////////////////

	cmd, err = ParseCommand("dir")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: "dir",
	})

	cmd, err = ParseCommand("dir", "--template=js")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: "dir",
	})

	cmd, err = ParseCommand("dir", "--template=jsx")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: "dir",
	})

	cmd, err = ParseCommand("dir", "--template=javascript")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: "dir",
	})

	cmd, err = ParseCommand("dir", "--template=ts")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: "dir",
	})

	cmd, err = ParseCommand("dir", "--template=tsx")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: "dir",
	})

	cmd, err = ParseCommand("dir", "--template=typescript")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: "dir",
	})
}
