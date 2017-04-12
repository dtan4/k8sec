package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeClient represents Kubernetes client and calculated namespace
type KubeClient struct {
	Clientset *kubernetes.Clientset
	Namespace string
}

// NewKubeClient creates new Kubernetes API client
func NewKubeClient(kubeconfig, context, namespace string) (*KubeClient, error) {
	if kubeconfig == "" {
		kubeconfig = clientcmd.RecommendedHomeFile
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: context},
	)

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return nil, err
	}

	ns := namespace

	if ns == "" {
		if rawConfig.Contexts[rawConfig.CurrentContext].Namespace == "" {
			ns = v1.NamespaceDefault
		} else {
			ns = rawConfig.Contexts[rawConfig.CurrentContext].Namespace
		}
	}

	return &KubeClient{
		Clientset: clientset,
		Namespace: namespace,
	}, nil
}

// GetSecret returns secret with the given name
func (c *KubeClient) GetSecret(name string) (*v1.Secret, error) {
	return c.Clientset.Core().Secrets(c.Namespace).Get(name)
}

// ListSecrets returns the list of Secrets
func (c *KubeClient) ListSecrets() (*v1.SecretList, error) {
	return c.Clientset.Core().Secrets(c.Namespace).List(v1.ListOptions{})
}
