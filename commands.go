package main

import (
	"github.com/dtan4/k8sec/command"
	"github.com/mitchellh/cli"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"set": func() (cli.Command, error) {
			return &command.SetCommand{
				Meta: *meta,
			}, nil
		},
		"unset": func() (cli.Command, error) {
			return &command.UnsetCommand{
				Meta: *meta,
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				Meta: *meta,
			}, nil
		},
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
				Revision: GitCommit,
				Name:     Name,
			}, nil
		},
	}
}
