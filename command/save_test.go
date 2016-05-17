package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestSaveCommand_implement(t *testing.T) {
	var _ cli.Command = &SaveCommand{}
}
