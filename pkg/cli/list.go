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
	"github.com/sapcc/gophercloud-limes/resources/v1/clusters"
	"github.com/sapcc/gophercloud-limes/resources/v1/domains"
	"github.com/sapcc/gophercloud-limes/resources/v1/projects"
)

// list retrieves information about all the clusters within the token scope.
func (c *Cluster) list() {
	provider, err := clientconfig.AuthenticatedClient(nil)
	if err != nil {
		kingpin.Fatalf("cannot connect to OpenStack: %v", err.Error())
	}
	limesClient, err := resources.NewLimesV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		kingpin.Fatalf("could not initialize Limes client: %v", err.Error())
	}

	c.IsList = true
	c.Result = clusters.List(limesClient, clusters.ListOpts{
		Area:     c.Opts.Area,
		Service:  c.Opts.Service,
		Resource: c.Opts.Resource,
	})
	if c.Result.Err != nil {
		kingpin.Fatalf("could not get clusters: %v", c.Result.Err)
	}
}

// list retrieves information about all the domains within the token scope.
func (d *Domain) list() {
	provider, err := clientconfig.AuthenticatedClient(nil)
	if err != nil {
		kingpin.Fatalf("cannot connect to OpenStack: %v", err)
	}
	identityClient, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		kingpin.Fatalf("could not initialize identity client: %v", err)
	}
	limesClient, err := resources.NewLimesV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		kingpin.Fatalf("could not initialize Limes client: %v", err)
	}

	nameOrID := d.ID
	count, err := d.find(identityClient)
	if err != nil {
		kingpin.Fatalf("could not find domain in token scope: %v", err)
	}
	if nameOrID == d.Name && count > 1 {
		kingpin.Fatalf("more than one domain exists with the name %v", d.Name)
	}

	d.IsList = true
	d.Result = domains.List(limesClient, domains.ListOpts{
		Area:     d.Opts.Area,
		Service:  d.Opts.Service,
		Resource: d.Opts.Resource,
	})
	if d.Result.Err != nil {
		kingpin.Fatalf("could not list domains: %v", d.Result.Err)
	}
}

// list retrieves information about all the projects within a specific domain.
func (p *Project) list() {
	provider, err := clientconfig.AuthenticatedClient(nil)
	if err != nil {
		kingpin.Fatalf("cannot connect to OpenStack: %v", err)
	}
	identityClient, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		kingpin.Fatalf("could not initialize identity client: %v", err)
	}
	limesClient, err := resources.NewLimesV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		kingpin.Fatalf("could not initialize Limes client: %v", err)
	}

	nameOrID := p.DomainID
	d := &Domain{ID: p.DomainID}
	count, err := d.find(identityClient)
	if err != nil {
		kingpin.Fatalf("could not find project's domain in token scope: %v", err)
	}
	if nameOrID == d.Name && count > 1 {
		kingpin.Fatalf("more than one domain exists with the name %v", d.Name)
	}

	p.DomainID = d.ID
	p.DomainName = d.Name

	p.IsList = true
	p.Result = projects.List(limesClient, p.DomainID, projects.ListOpts{
		Area:     p.Opts.Area,
		Service:  p.Opts.Service,
		Resource: p.Opts.Resource,
	})
	if p.Result.Err != nil {
		kingpin.Fatalf("could not list projects: %v", d.Result.Err)
	}
}
