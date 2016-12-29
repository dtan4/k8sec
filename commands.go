package main

import (
	"github.com/dtan4/k8sec/command"
	"github.com/mitchellh/cli"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"load": func() (cli.Command, error) {
			return &command.LoadCommand{
				Meta: *meta,
			}, nil
		},
		"save": func() (cli.Command, error) {
			return &command.SaveCommand{
				Meta: *meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  Version,
				Revision: Revision,
				Name:     Name,
			}, nil
		},
	}
}
