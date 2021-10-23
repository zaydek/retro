package main

import (
	"fmt"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	messages := api.FormatMessages([]api.Message{
		{
			Text: // "An uncaught runtime error occurred" +
			// "\n\n" +
			`Error: Oops
    at http://localhost:8000/client.js:896:9
    at http://localhost:8000/client.js:949:3`,
		},
	}, api.FormatMessagesOptions{
		Color: true,
		// Kind:       api.WarningMessage,
		Kind:          api.ErrorMessage,
		TerminalWidth: 80,
	})
	str := strings.TrimRight(strings.Join(messages, ""), "\n")
	fmt.Println(str)
}
