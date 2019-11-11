package auth

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/utils/client"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sapcc/gophercloud-limes/resources"
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

	limesClient, err := resources.NewLimesV1(provider, gophercloud.EndpointOpts{})
	errors.Handle(err, "could not initialize Limes client")

	return identityClient, limesClient
}
