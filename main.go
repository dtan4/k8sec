package main

import (
	"os"

	"github.com/dtan4/k8sec/cmd"
)

func main() {
	cmd.Execute(os.Stdin, os.Stdout, os.Args)
}
