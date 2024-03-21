// Code generated by codegen. DO NOT EDIT.

package upgrade

import (
	v1 "github.com/buroa/system-upgrade-controller/pkg/generated/controllers/upgrade.cattle.io/v1"
	"github.com/rancher/lasso/pkg/controller"
)

type Interface interface {
	V1() v1.Interface
}

type group struct {
	controllerFactory controller.SharedControllerFactory
}

// New returns a new Interface.
func New(controllerFactory controller.SharedControllerFactory) Interface {
	return &group{
		controllerFactory: controllerFactory,
	}
}

func (g *group) V1() v1.Interface {
	return v1.New(g.controllerFactory)
}
