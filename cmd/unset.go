package cmd

import (
	"fmt"

	"github.com/wantedly/k8sec/k8s"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset KEY1 [KEY2 ...]",
	Short: "Unset secrets",
	RunE:  doUnset,
}

func doUnset(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("Too few arguments.")
	}
	name := args[0]

	k8sclient, err := k8s.NewKubeClient(rootOpts.kubeconfig, rootOpts.context)
	if err != nil {
		return errors.Wrap(err, "Failed to initialize Kubernetes API client.")
	}

	var namespace string

	if rootOpts.namespace != "" {
		namespace = rootOpts.namespace
	} else {
		namespace = k8sclient.DefaultNamespace()
	}

	s, err := k8sclient.GetSecret(namespace, name)
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

	_, err = k8sclient.UpdateSecret(namespace, s)
	if err != nil {
		return errors.Wrapf(err, "Failed to unset secret. name=%s", name)
	}

	fmt.Println(s.Name)

	return nil
}

func init() {
	RootCmd.AddCommand(unsetCmd)
}
