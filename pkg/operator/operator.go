/*
Copyright © 2022 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package operator

import (
	"context"

	v1 "github.com/rancher-sandbox/rancheros-operator/pkg/apis/rancheros.cattle.io/v1"
	"github.com/rancher-sandbox/rancheros-operator/pkg/clients"
	"github.com/rancher-sandbox/rancheros-operator/pkg/controllers/inventory"
	"github.com/rancher-sandbox/rancheros-operator/pkg/controllers/managedos"
	"github.com/rancher-sandbox/rancheros-operator/pkg/controllers/managedosversionchannel"

	"github.com/rancher-sandbox/rancheros-operator/pkg/controllers/registration"
	"github.com/rancher-sandbox/rancheros-operator/pkg/server"
	"github.com/rancher/steve/pkg/aggregation"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func Run(ctx context.Context, settings ...Setting) error {
	o := &options{}

	if err := o.apply(settings...); err != nil {
		return err
	}

	restConfig, err := config.GetConfig()
	if err != nil {
		logrus.Fatalf("failed to find kubeconfig: %v", err)
	}

	clients, err := clients.NewFromConfig(restConfig)
	if err != nil {
		logrus.Fatalf("Error building controller: %s", err.Error())
	}

	factory, err := crd.NewFactoryFromClient(restConfig)
	if err != nil {
		logrus.Fatalf("Failed to create CRD factory: %v", err)
	}

	err = factory.BatchCreateCRDs(ctx,
		crd.CRD{
			SchemaObject: v1.ManagedOSImage{},
			Status:       true,
		},
		crd.CRD{
			SchemaObject: v1.MachineInventory{},
			Status:       true,
		},
		crd.CRD{
			SchemaObject: v1.MachineRegistration{},
			Status:       true,
		},
		crd.CRD{
			SchemaObject: v1.ManagedOSVersion{},
			Status:       true,
		},
		crd.CRD{
			SchemaObject: v1.ManagedOSVersionChannel{},
			Status:       true,
		},
	).BatchWait()
	if err != nil {
		logrus.Fatalf("Failed to create CRDs: %v", err)
	}

	managedos.Register(ctx, clients)
	inventory.Register(ctx, clients)
	registration.Register(ctx, clients)
	managedosversionchannel.Register(ctx, o.requeuer, clients)

	aggregation.Watch(ctx, clients.Core.Secret(), o.namespace, "rancheros-operator", server.New(clients))

	for _, s := range o.services {
		go s(ctx, clients) //nolint:golint,errcheck
	}

	return clients.Start(ctx)
}
