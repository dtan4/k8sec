package command

import (
	"strings"
)

type SaveCommand struct {
	Meta
}

func (c *SaveCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *SaveCommand) Synopsis() string {
	return ""
}

func (c *SaveCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
