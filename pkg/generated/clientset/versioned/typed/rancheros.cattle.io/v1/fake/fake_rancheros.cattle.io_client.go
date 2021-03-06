/*
Copyright 2022 SUSE LLC

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

// Code generated by main. DO NOT EDIT.

package fake

import (
	v1 "github.com/rancher-sandbox/rancheros-operator/pkg/generated/clientset/versioned/typed/rancheros.cattle.io/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeRancherosV1 struct {
	*testing.Fake
}

func (c *FakeRancherosV1) MachineInventories(namespace string) v1.MachineInventoryInterface {
	return &FakeMachineInventories{c, namespace}
}

func (c *FakeRancherosV1) MachineRegistrations(namespace string) v1.MachineRegistrationInterface {
	return &FakeMachineRegistrations{c, namespace}
}

func (c *FakeRancherosV1) ManagedOSImages(namespace string) v1.ManagedOSImageInterface {
	return &FakeManagedOSImages{c, namespace}
}

func (c *FakeRancherosV1) ManagedOSVersions(namespace string) v1.ManagedOSVersionInterface {
	return &FakeManagedOSVersions{c, namespace}
}

func (c *FakeRancherosV1) ManagedOSVersionChannels(namespace string) v1.ManagedOSVersionChannelInterface {
	return &FakeManagedOSVersionChannels{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeRancherosV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
