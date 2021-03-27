package create_retro_app

import (
	"testing"

	"github.com/zaydek/retro/pkg/expect"
)

func TestParseArguments(t *testing.T) {
	var cmd Command

	cmd = parseArguments(".")
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: ".",
	})

	cmd = parseArguments("--template=javascript", ".")
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: ".",
	})

	cmd = parseArguments("--template=typescript", ".")
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: ".",
	})

	cmd = parseArguments("app-name")
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: "app-name",
	})

	cmd = parseArguments("--template=javascript", "app-name")
	expect.DeepEqual(t, cmd, Command{
		Template:  "javascript",
		Directory: "app-name",
	})

	cmd = parseArguments("--template=typescript", "app-name")
	expect.DeepEqual(t, cmd, Command{
		Template:  "typescript",
		Directory: "app-name",
	})
}
