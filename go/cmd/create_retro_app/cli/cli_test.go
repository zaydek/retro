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
		Template:  "starter",
		Directory: ".",
	})

	cmd, err = ParseCommand(".", "--template=starter")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "starter",
		Directory: ".",
	})

	cmd, err = ParseCommand(".", "--template=sass")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "sass",
		Directory: ".",
	})

	cmd, err = ParseCommand(".", "--template=mdx")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "mdx",
		Directory: ".",
	})

	//////////////////////////////////////////////////////////////////////////////

	cmd, err = ParseCommand("app")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "starter",
		Directory: "app",
	})

	cmd, err = ParseCommand("app", "--template=starter")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "starter",
		Directory: "app",
	})

	cmd, err = ParseCommand("app", "--template=sass")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "sass",
		Directory: "app",
	})

	cmd, err = ParseCommand("app", "--template=mdx")
	check(t, err)
	expect.DeepEqual(t, cmd, Command{
		Template:  "mdx",
		Directory: "app",
	})
}
