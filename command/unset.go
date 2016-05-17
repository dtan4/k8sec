package command

import (
	"strings"
)

type UnsetCommand struct {
	Meta
}

func (c *UnsetCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *UnsetCommand) Synopsis() string {
	return ""
}

func (c *UnsetCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
