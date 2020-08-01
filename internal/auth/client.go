// Copyright 2020 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/utils/client"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sapcc/gophercloud-sapcc/clients"
	"github.com/sapcc/limesctl/internal/errors"
)

// ServiceClients authenticates against OpenStack and returns respective ServiceClients
// for Keystone and Limes.
func ServiceClients(debug bool) (identityV3, limesV1 *gophercloud.ServiceClient) {
	ao, err := clientconfig.AuthOptions(nil)
	errors.Handle(err, "can not get auth variables")

	provider, err := openstack.NewClient(ao.IdentityEndpoint)
	errors.Handle(err, "can not create an OpenStack client")

	if debug {
		provider.HTTPClient = http.Client{
			Transport: &client.RoundTripper{
				Rt:     &http.Transport{},
				Logger: &client.DefaultLogger{},
			},
		}
	}

	err = openstack.Authenticate(provider, *ao)
	errors.Handle(err, "can not connect to OpenStack")

	identityClient, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	errors.Handle(err, "could not initialize identity client")

	limesClient, err := clients.NewLimesV1(provider, gophercloud.EndpointOpts{})
	errors.Handle(err, "could not initialize Limes client")

	return identityClient, limesClient
}
