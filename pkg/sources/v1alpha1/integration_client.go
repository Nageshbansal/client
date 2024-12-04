package v1alpha1

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/client/pkg/context"
	knerrors "knative.dev/client/pkg/errors"
	"knative.dev/client/pkg/util"
	sourcesv1alpha1 "knative.dev/eventing/pkg/apis/sources/v1alpha1"
	"knative.dev/eventing/pkg/client/clientset/versioned/scheme"
	clientv1alpha1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/sources/v1alpha1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type KnIntegrationSourcesClient interface {
	GetIntegrationSource(ctx context.Context, name string) (*sourcesv1alpha1.IntegrationSource, error)

	CreateIntegrationSource(ctx context.Context, source *sourcesv1alpha1.IntegrationSource) error

	UpdateIntegrationSource(ctx context.Context, source *sourcesv1alpha1.IntegrationSource) error

	DeleteIntegrationSource(ctx context.Context, name string) error

	ListIntegrationSource(ctx context.Context) (*sourcesv1alpha1.IntegrationSourceList, error)

	Namespace() string
}

type integrationSourcesClient struct {
	client    clientv1alpha1.IntegrationSourceInterface
	namespace string
}

func newKnIntegrationSourcesClient(client clientv1alpha1.IntegrationSourceInterface, namespace string) KnIntegrationSourcesClient {
	return &integrationSourcesClient{
		client:    client,
		namespace: namespace,
	}
}

func (c *integrationSourcesClient) Namespace() string {
	return c.namespace
}

// GetIntegrationSource returns the available IntegrationSource
func (c *integrationSourcesClient) GetIntegrationSource(ctx context.Context, name string) (*sourcesv1alpha1.IntegrationSource, error) {
	source, err := c.client.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, knerrors.GetError(err)
	}
	err = updateIntegrationSourceGVK(source)
	if err != nil {
		return nil, knerrors.GetError(err)
	}
	return source, nil
}

func updateIntegrationSourceGVK(obj runtime.Object) error {
	return util.UpdateGroupVersionKindWithScheme(obj, sourcesv1alpha1.SchemeGroupVersion, scheme.Scheme)
}

func updateIntegrationSourceListGVK(sourcesList *sourcesv1alpha1.IntegrationSourceList) (*sourcesv1alpha1.IntegrationSourceList, error) {
	sourcesListNew := sourcesList.DeepCopy()
	err := updateIntegrationSourceGVK(sourcesListNew)
	if err != nil {
		return nil, knerrors.GetError(err)
	}

	sourcesListNew.Items = make([]sourcesv1alpha1.IntegrationSource, len(sourcesList.Items))
	for idx, source := range sourcesList.Items {
		sourceClone := source.DeepCopy()
		err := updateIntegrationSourceGVK(sourceClone)
		if err != nil {
			return nil, knerrors.GetError(err)
		}
		sourcesListNew.Items[idx] = *sourceClone
	}
	return sourcesListNew, nil
}

func (c *integrationSourcesClient) CreateIntegrationSource(ctx context.Context, integrationSource *sourcesv1alpha1.IntegrationSource) error {
	if integrationSource.Spec.Sink.Ref == nil && integrationSource.Spec.Sink.URI == nil {
		return fmt.Errorf("a sink is required for creating an IntegrationSource")
	}
	_, err := c.client.Create(ctx, integrationSource, metav1.CreateOptions{})
	if err != nil {
		return knerrors.GetError(err)
	}
	return nil
}

func (c *integrationSourcesClient) UpdateIntegrationSource(ctx context.Context, integrationSource *sourcesv1alpha1.IntegrationSource) error {
	_, err := c.client.Update(ctx, integrationSource, metav1.UpdateOptions{})
	if err != nil {
		return knerrors.GetError(err)
	}
	return nil
}

func (c *integrationSourcesClient) DeleteIntegrationSource(ctx context.Context, name string) error {
	err := c.client.Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return knerrors.GetError(err)
	}
	return nil
}

func (c *integrationSourcesClient) ListIntegrationSource(ctx context.Context) (*sourcesv1alpha1.IntegrationSourceList, error) {

	sourcesList, err := c.client.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, knerrors.GetError(err)
	}
	return updateIntegrationSourceListGVK(sourcesList)
}

// IntegrationSource Builder

type IntegrationSourceBuilder struct {
	integrationSource *sourcesv1alpha1.IntegrationSource
}

func NewIntegrationSourceBuilder(name string) *IntegrationSourceBuilder {
	return &IntegrationSourceBuilder{integrationSource: &sourcesv1alpha1.IntegrationSource{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}}
}
func NewIntegrationSourceFromExisting(i *sourcesv1alpha1.IntegrationSource) *IntegrationSourceBuilder {
	return &IntegrationSourceBuilder{
		integrationSource: i.DeepCopy(),
	}
}

func (i *IntegrationSourceBuilder) Sink(sink duckv1.Destination) *IntegrationSourceBuilder {
	i.integrationSource.Spec.Sink = sink
	return i
}

func (i *IntegrationSourceBuilder) AWS(aws sourcesv1alpha1.Aws) *IntegrationSourceBuilder {
	i.integrationSource.Spec.Aws = aws.DeepCopy()
	return i
}

func (i *IntegrationSourceBuilder) Timer(timer sourcesv1alpha1.Timer) *IntegrationSourceBuilder {
	i.integrationSource.Spec.Timer = timer.DeepCopy()
	return i
}

func (i *IntegrationSourceBuilder) Build() *sourcesv1alpha1.IntegrationSource {
	return i.integrationSource
}
