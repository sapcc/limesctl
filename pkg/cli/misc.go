/*******************************************************************************
*
* Copyright 2018 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package cli

import (
	"github.com/alecthomas/kingpin"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sapcc/gophercloud-limes/resources"
)

// getServiceClients authenticates against OpenStack and returns respective ServiceClients
// for Keystone and Limes.
func getServiceClients() (identityV3, limesV1 *gophercloud.ServiceClient) {
	provider, err := clientconfig.AuthenticatedClient(nil)
	handleError("can not connect to OpenStack", err)

	identityClient, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	handleError("could not initialize identity client", err)

	limesClient, err := resources.NewLimesV1(provider, gophercloud.EndpointOpts{})
	handleError("could not initialize Limes client", err)

	return identityClient, limesClient
}

// handleError is a convenient wrapper around kingpin.Fatalf.
func handleError(str string, err error) {
	if err == nil {
		return
	}

	kingpin.Fatalf("%v: %v", str, err)
}
