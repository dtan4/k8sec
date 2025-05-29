package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/dtan4/k8sec/pkg/client"
	"github.com/spf13/cobra"
)

func newUnsetCmd(out io.Writer) *cobra.Command {
	unsetCmd := &cobra.Command{
		Use:   "unset KEY1 [KEY2 ...]",
		Short: "Unset secrets",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("too few arguments")
			}

			ctx := context.Background()

			k8sclient, err := client.New(rootOpts.kubeconfig, rootOpts.context)
			if err != nil {
				return fmt.Errorf("initialize Kubernetes API client: %w", err)
			}

			var namespace string

			if rootOpts.namespace != "" {
				namespace = rootOpts.namespace
			} else {
				namespace = k8sclient.DefaultNamespace()
			}

			return runUnset(ctx, k8sclient, namespace, args, out)
		},
	}

	return unsetCmd
}

func runUnset(ctx context.Context, k8sclient client.Client, namespace string, args []string, out io.Writer) error {
	name := args[0]

	s, err := k8sclient.GetSecret(ctx, namespace, name)
	if err != nil {
		return fmt.Errorf("get current secret %q: %w", name, err)
	}

	for _, k := range args[1:] {
		_, ok := s.Data[k]
		if !ok {
			return fmt.Errorf("the key %s does not exist", k)
		}

		delete(s.Data, k)
	}

	_, err = k8sclient.UpdateSecret(ctx, namespace, s)
	if err != nil {
		return fmt.Errorf("unset secret %q: %w", name, err)
	}

	fmt.Fprintln(out, s.Name)

	return nil
}

func init() {
}
