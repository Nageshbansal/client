package integration

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"knative.dev/client/pkg/commands"
	"knative.dev/client/pkg/sources/v1alpha1"
	sourcev1alpha1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/sources/v1alpha1"
)

var integrationSourceClientFactory func(config clientcmd.ClientConfig, namespace string) (v1alpha1.KnIntegrationSourcesClient, error)

func newIntegrationSourceClient(p *commands.KnParams, cmd *cobra.Command) (v1alpha1.KnIntegrationSourcesClient, error) {

	namespace, err := p.GetNamespace(cmd)
	if err != nil {
		return nil, err
	}

	if integrationSourceClientFactory != nil {
		config, err := p.GetClientConfig()
		if err != nil {
			return nil, err
		}
		return integrationSourceClientFactory(config, namespace)
	}

	clientConfig, err := p.RestConfig()
	if err != nil {
		return nil, err
	}
	client, err := sourcev1alpha1.NewForConfig(clientConfig)
	return v1alpha1.NewKnSourcesClient(client, namespace).IntegrationSourcesClient(), nil
}
