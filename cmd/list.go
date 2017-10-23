package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/k8sec/k8s"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var listOpts = struct {
	base64encode bool
}{}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets",
	Long: `List secrets

$ k8sec list rails
NAME    TYPE    KEY             VALUE
rails   Opaque  database-url    "postgres://example.com:5432/dbname"

Show values as base64-encoded string:

$ k8sec list --base64 rails
NAME    TYPE    KEY             VALUE
rails   Opaque  database-url    cG9zdGdyZXM6Ly9leGFtcGxlLmNvbTo1NDMyL2RibmFtZQ==
`,
	RunE: doList,
}

type Secret struct {
	Name  string
	Type  string
	Key   string
	Value string
}

func doList(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("Too many arguments.")
	}

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

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, strings.Join([]string{"NAME", "TYPE", "KEY", "VALUE"}, "\t"))

	var v string

	secrets := []Secret{}

	if len(args) == 1 {
		secret, err := k8sclient.GetSecret(namespace, args[0])
		if err != nil {
			return errors.Wrap(err, "Failed to retrieve secrets.")
		}

		for key, value := range secret.Data {
			if listOpts.base64encode {
				v = base64.StdEncoding.EncodeToString(value)
			} else {
				v = strconv.Quote(string(value))
			}

			secrets = append(secrets, Secret{
				Name:  secret.Name,
				Type:  string(secret.Type),
				Key:   key,
				Value: v,
			})
		}

		// sort by KEY
		sort.Slice(secrets, func(i, j int) bool {
			return secrets[i].Key < secrets[j].Key
		})
	} else {
		ss, err := k8sclient.ListSecrets(namespace)
		if err != nil {
			return errors.Wrap(err, "Failed to retrieve secrets.")
		}

		for _, secret := range ss.Items {
			kvs := []struct {
				k, v string
			}{}

			for key, value := range secret.Data {
				if listOpts.base64encode {
					v = base64.StdEncoding.EncodeToString(value)
				} else {
					v = strconv.Quote(string(value))
				}

				kvs = append(kvs, struct {
					k, v string
				}{
					k: key,
					v: v,
				})
			}

			// sort by KEY
			sort.Slice(kvs, func(i, j int) bool {
				return kvs[i].k < kvs[j].k
			})

			for _, kv := range kvs {
				secrets = append(secrets, Secret{
					Name:  secret.Name,
					Type:  string(secret.Type),
					Key:   kv.k,
					Value: kv.v,
				})
			}
		}
	}

	for _, secret := range secrets {
		fmt.Fprintln(w, strings.Join([]string{secret.Name, secret.Type, secret.Key, secret.Value}, "\t"))
	}

	w.Flush()

	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&listOpts.base64encode, "base64", false, "Show values as base64-encoded string")
}
