package command

import (
	"strings"
)

type LoadCommand struct {
	Meta
}

func (c *LoadCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *LoadCommand) Synopsis() string {
	return ""
}

func (c *LoadCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
