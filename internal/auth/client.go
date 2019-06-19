package auth

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sapcc/gophercloud-limes/resources"
	"github.com/sapcc/limesctl/internal/errors"
)

// ServiceClients authenticates against OpenStack and returns respective ServiceClients
// for Keystone and Limes.
func ServiceClients() (identityV3, limesV1 *gophercloud.ServiceClient) {
	provider, err := clientconfig.AuthenticatedClient(nil)
	errors.Handle(err, "can not connect to OpenStack")

	identityClient, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	errors.Handle(err, "could not initialize identity client")

	limesClient, err := resources.NewLimesV1(provider, gophercloud.EndpointOpts{})
	errors.Handle(err, "could not initialize Limes client")

	return identityClient, limesClient
}
