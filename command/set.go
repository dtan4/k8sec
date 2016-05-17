package command

import (
	"strings"
)

type SetCommand struct {
	Meta
}

func (c *SetCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *SetCommand) Synopsis() string {
	return ""
}

func (c *SetCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
