package cmd

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

type fakeClient struct {
	defaultNamespace     string
	getSecretResponse    *v1.Secret
	listSecretsResponse  *v1.SecretList
	updateSecretResponse *v1.Secret
	err                  error
}

func (c *fakeClient) DefaultNamespace() string {
	return c.defaultNamespace
}

func (c *fakeClient) CreateSecret(ctx context.Context, namespace string, secret *v1.Secret) (*v1.Secret, error) {
	return secret, c.err
}

func (c *fakeClient) GetSecret(ctx context.Context, namespace, name string) (*v1.Secret, error) {
	return c.getSecretResponse, c.err
}

func (c *fakeClient) ListSecrets(ctx context.Context, namespace string) (*v1.SecretList, error) {
	return c.listSecretsResponse, c.err
}

func (c *fakeClient) UpdateSecret(ctx context.Context, namespace string, secret *v1.Secret) (*v1.Secret, error) {
	return c.updateSecretResponse, c.err
}
