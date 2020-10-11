package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/dtan4/k8sec/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newUnsetCmd(out io.Writer) *cobra.Command {
	unsetCmd := &cobra.Command{
		Use:   "unset KEY1 [KEY2 ...]",
		Short: "Unset secrets",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("Too few arguments.")
			}

			ctx := context.Background()

			k8sclient, err := client.New(rootOpts.kubeconfig, rootOpts.context)
			if err != nil {
				return errors.Wrap(err, "Failed to initialize Kubernetes API client.")
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
		return errors.Wrapf(err, "Failed to get current secret. name=%s", name)
	}

	for _, k := range args[1:] {
		_, ok := s.Data[k]
		if !ok {
			return errors.Errorf("The key %s does not exist.", k)
		}

		delete(s.Data, k)
	}

	_, err = k8sclient.UpdateSecret(ctx, namespace, s)
	if err != nil {
		return errors.Wrapf(err, "Failed to unset secret. name=%s", name)
	}

	fmt.Fprintln(out, s.Name)

	return nil
}

func init() {
}
