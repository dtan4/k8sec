package command

import (
	"fmt"
	"os"
	"strings"
)

type SaveCommand struct {
	Meta
}

func (c *SaveCommand) Run(args []string) int {
	fmt.Fprintln(os.Stderr, "save command is currently not implemented.")

	return 1
}

func (c *SaveCommand) Synopsis() string {
	return ""
}

func (c *SaveCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
