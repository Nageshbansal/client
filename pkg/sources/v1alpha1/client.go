package v1alpha1

import (
	clientv1alpha1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/sources/v1alpha1"
)

// KnSourcesClient to Eventing Sources. All methods are relative to the
// namespace specified during construction
type KnSourcesClient interface {
	// IntegrationSourcesClient , Get client for Integration sources
	IntegrationSourcesClient() KnIntegrationSourcesClient
}

// sourcesClient is a combination of Sources client interface and namespace
// Temporarily help to add sources dependencies
// May be changed when adding real sources features
type sourcesClient struct {
	client    clientv1alpha1.SourcesV1alpha1Interface
	namespace string
}

// NewKnSourcesClient for managing all eventing built-in sources
func NewKnSourcesClient(client clientv1alpha1.SourcesV1alpha1Interface, namespace string) KnSourcesClient {
	return &sourcesClient{
		client:    client,
		namespace: namespace,
	}
}

// IntegrationSourcesClient get the client for dealing with Integration sources
func (c *sourcesClient) IntegrationSourcesClient() KnIntegrationSourcesClient {
	return newKnIntegrationSourcesClient(c.client.IntegrationSources(c.namespace), c.namespace)
}
