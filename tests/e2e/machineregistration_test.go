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

package e2e_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	http "github.com/rancher-sandbox/ele-testhelpers/http"
	kubectl "github.com/rancher-sandbox/ele-testhelpers/kubectl"

	"github.com/rancher-sandbox/rancheros-operator/tests/catalog"
)

var _ = Describe("MachineRegistration e2e tests", func() {
	k := kubectl.New()
	Context("registration", func() {

		AfterEach(func() {
			kubectl.New().Delete("machineregistration", "-n", "fleet-default", "machine-registration")
		})

		It("creates a machine registration resource and a URL attaching CA certificate", func() {
			mr := catalog.NewMachineRegistration("machine-registration", map[string]interface{}{
				"install":   map[string]string{"device": "/dev/vda"},
				"rancheros": map[string]interface{}{"install": map[string]string{"isoUrl": "https://something.example.com"}},
				"users": []map[string]string{
					{
						"name":   "root",
						"passwd": "root",
					},
				},
			})

			Eventually(func() error {
				return k.ApplyYAML("fleet-default", "machine-registration", mr)
			}, 2*time.Minute, 2*time.Second).ShouldNot(HaveOccurred())

			var url string
			Eventually(func() string {
				e, err := kubectl.GetData("fleet-default", "machineregistration", "machine-registration", `jsonpath={.status.registrationURL}`)
				if err != nil {
					fmt.Println(err)
				}
				url = string(e)
				return string(e)
			}, 1*time.Minute, 2*time.Second).Should(
				And(
					ContainSubstring(fmt.Sprintf("%s.%s/v1-rancheros/registration", externalIP, magicDNS)),
				),
			)

			out, err := http.GetInsecure(fmt.Sprintf("https://%s", url))
			Expect(err).ToNot(HaveOccurred())

			Expect(out).Should(
				And(
					ContainSubstring("BEGIN CERTIFICATE"),
					ContainSubstring(fmt.Sprintf("%s.%s/v1-rancheros/registration", externalIP, magicDNS)),
				),
			)
		})
	})
})
