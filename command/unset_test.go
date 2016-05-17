package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestUnsetCommand_implement(t *testing.T) {
	var _ cli.Command = &UnsetCommand{}
}
