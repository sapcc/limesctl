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

package core

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sapcc/gophercloud-limes/resources"
	"github.com/sapcc/limesctl/internal/errors"
)

// getServiceClients authenticates against OpenStack and returns respective ServiceClients
// for Keystone and Limes.
func getServiceClients() (identityV3, limesV1 *gophercloud.ServiceClient) {
	provider, err := clientconfig.AuthenticatedClient(nil)
	errors.Handle(err, "can not connect to OpenStack")

	identityClient, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	errors.Handle(err, "could not initialize identity client")

	limesClient, err := resources.NewLimesV1(provider, gophercloud.EndpointOpts{})
	errors.Handle(err, "could not initialize Limes client")

	return identityClient, limesClient
}
