package command

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dtan4/k8sec/k8s"
)

type SetCommand struct {
	Meta
}

func (c *SetCommand) Run(args []string) int {
	var (
		arguments  []string
		kubeconfig string
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")

	if err := flags.Parse(args[0:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}

	if len(arguments) < 2 {
		fmt.Fprintln(os.Stderr, "Too few arguments. Example: $ k8sec set rails RAILS_ENV=production")
		return 1
	}

	_ = arguments[0]

	secrets := map[string]string{}

	for _, kv := range arguments[1:] {
		ary := strings.SplitN(kv, "=", 2)

		if len(ary) != 2 {
			fmt.Fprintln(os.Stderr, "Argument should be key=value format. Given argument: "+kv)
			return 1
		}

		k, v := ary[0], ary[1]
		secrets[k] = v
	}

	_, err := k8s.NewKubeClient(kubeconfig)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	for k, v := range secrets {
		fmt.Println(k + ": " + v)
	}

	return 0
}

func (c *SetCommand) Synopsis() string {
	return ""
}

func (c *SetCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
