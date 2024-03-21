// Code generated by codegen. DO NOT EDIT.

package upgrade

import (
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/v2/pkg/generic"
	"k8s.io/client-go/rest"
)

type Factory struct {
	*generic.Factory
}

func NewFactoryFromConfigOrDie(config *rest.Config) *Factory {
	f, err := NewFactoryFromConfig(config)
	if err != nil {
		panic(err)
	}
	return f
}

func NewFactoryFromConfig(config *rest.Config) (*Factory, error) {
	return NewFactoryFromConfigWithOptions(config, nil)
}

func NewFactoryFromConfigWithNamespace(config *rest.Config, namespace string) (*Factory, error) {
	return NewFactoryFromConfigWithOptions(config, &FactoryOptions{
		Namespace: namespace,
	})
}

type FactoryOptions = generic.FactoryOptions

func NewFactoryFromConfigWithOptions(config *rest.Config, opts *FactoryOptions) (*Factory, error) {
	f, err := generic.NewFactoryFromConfigWithOptions(config, opts)
	return &Factory{
		Factory: f,
	}, err
}

func NewFactoryFromConfigWithOptionsOrDie(config *rest.Config, opts *FactoryOptions) *Factory {
	f, err := NewFactoryFromConfigWithOptions(config, opts)
	if err != nil {
		panic(err)
	}
	return f
}

func (c *Factory) Upgrade() Interface {
	return New(c.ControllerFactory())
}

func (c *Factory) WithAgent(userAgent string) Interface {
	return New(controller.NewSharedControllerFactoryWithAgent(userAgent, c.ControllerFactory()))
}
