package create_retro_app

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/zaydek/retro/pkg/loggers"
)

func parseArguments(arguments ...string) Command {
	flagset := flag.NewFlagSet("", flag.ContinueOnError)
	flagset.SetOutput(ioutil.Discard)

	cmd := Command{}
	flagset.StringVar(&cmd.Template, "template", "javascript", "")
	if err := flagset.Parse(arguments); err != nil {
		fmt.Println(usage)
		os.Exit(2)
	}
	if cmd.Template != "javascript" && cmd.Template != "typescript" {
		loggers.Error("--template must be javascript or typescript. " +
			"Here’s what you can do:\n\n" +
			"- create-retro-app --template=javascript app-name\n\n" +
			"Or\n\n" +
			"- create-retro-app --template=typescript app-name")
		os.Exit(2)
	}
	if len(flagset.Args()) == 0 {
		loggers.Error("It looks like you’re trying to run create-retro-app in the current directory. " +
			"In that case, use '.' explicitly.\n\n" +
			"- create-retro-app .\n\n" +
			"Or\n\n" +
			"- create-retro-app app-name")
		os.Exit(2)
	}
	cmd.Directory = "."
	if len(flagset.Args()) > 0 {
		cmd.Directory = flagset.Args()[0]
	}
	return cmd
}
