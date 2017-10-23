package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// KubeClient represents Kubernetes client and calculated namespace
type KubeClient struct {
	clientset *kubernetes.Clientset
	rawConfig api.Config
}

// NewKubeClient creates new Kubernetes API client
func NewKubeClient(kubeconfig, context string) (*KubeClient, error) {
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

	return &KubeClient{
		clientset: clientset,
		rawConfig: rawConfig,
	}, nil
}

// DefaultNamespace returns the default namespace in kubeconfig
func (c *KubeClient) DefaultNamespace() string {
	var namespace string

	if c.rawConfig.Contexts[c.rawConfig.CurrentContext].Namespace == "" {
		namespace = v1.NamespaceDefault
	} else {
		namespace = c.rawConfig.Contexts[c.rawConfig.CurrentContext].Namespace
	}

	return namespace
}

// CreateSecret creates new Secret
func (c *KubeClient) CreateSecret(namespace string, secret *v1.Secret) (*v1.Secret, error) {
	return c.clientset.Core().Secrets(namespace).Create(secret)
}

// GetSecret returns secret with the given name
func (c *KubeClient) GetSecret(namespace, name string) (*v1.Secret, error) {
	return c.clientset.Core().Secrets(namespace).Get(name)
}

// ListSecrets returns the list of Secrets
func (c *KubeClient) ListSecrets(namespace string) (*v1.SecretList, error) {
	return c.clientset.Core().Secrets(namespace).List(v1.ListOptions{})
}

// UpdateSecret updates the existed secret
func (c *KubeClient) UpdateSecret(namespace string, secret *v1.Secret) (*v1.Secret, error) {
	return c.clientset.Core().Secrets(namespace).Update(secret)
}
