package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestSetCommand_implement(t *testing.T) {
	var _ cli.Command = &SetCommand{}
}
