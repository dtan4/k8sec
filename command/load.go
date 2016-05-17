package command

import (
	"fmt"
	"os"
	"strings"
)

type LoadCommand struct {
	Meta
}

func (c *LoadCommand) Run(args []string) int {
	fmt.Fprintln(os.Stderr, "load command is currently not implemented.")

	return 1
}

func (c *LoadCommand) Synopsis() string {
	return ""
}

func (c *LoadCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
