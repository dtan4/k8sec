package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestLoadCommand_implement(t *testing.T) {
	var _ cli.Command = &LoadCommand{}
}
