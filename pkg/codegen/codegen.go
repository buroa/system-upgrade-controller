package main

import (
	v1 "github.com/buroa/system-upgrade-controller/pkg/apis/upgrade.cattle.io/v1"
	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	controllergen.Run(args.Options{
		Boilerplate:   "/dev/null",
		OutputPackage: "github.com/buroa/system-upgrade-controller/pkg/generated",
		Groups: map[string]args.Group{
			"upgrade.cattle.io": {
				Types: []interface{}{
					v1.Plan{},
				},
				GenerateTypes:   true,
				GenerateClients: true,
			},
		},
	})
}
