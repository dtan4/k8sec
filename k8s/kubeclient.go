package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeClient represents Kubernetes client and calculated namespace
type KubeClient struct {
	clientset *kubernetes.Clientset
	namespace string
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
		clientset: clientset,
		namespace: namespace,
	}, nil
}

// CreateSecret creates new Secret
func (c *KubeClient) CreateSecret(secret *v1.Secret) (*v1.Secret, error) {
	return c.clientset.Core().Secrets(c.namespace).Create(secret)
}

// GetSecret returns secret with the given name
func (c *KubeClient) GetSecret(name string) (*v1.Secret, error) {
	return c.clientset.Core().Secrets(c.namespace).Get(name)
}

// ListSecrets returns the list of Secrets
func (c *KubeClient) ListSecrets() (*v1.SecretList, error) {
	return c.clientset.Core().Secrets(c.namespace).List(v1.ListOptions{})
}

// UpdateSecret updates the existed secret
func (c *KubeClient) UpdateSecret(secret *v1.Secret) (*v1.Secret, error) {
	return c.clientset.Core().Secrets(c.namespace).Update(secret)
}
