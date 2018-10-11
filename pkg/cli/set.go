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
	idD "github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sapcc/gophercloud-limes/resources"
	"github.com/sapcc/gophercloud-limes/resources/v1/clusters"
	"github.com/sapcc/gophercloud-limes/resources/v1/domains"
	"github.com/sapcc/gophercloud-limes/resources/v1/projects"
	"github.com/sapcc/limes/pkg/api"
)

// set updates the resource capacities for a cluster.
func (c *Cluster) set(q interface{}) {
	provider, err := clientconfig.AuthenticatedClient(nil)
	if err != nil {
		kingpin.Fatalf("can not connect to OpenStack: %v", err)
	}
	limesClient, err := resources.NewLimesV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		kingpin.Fatalf("could not initialize Limes client: %v", err)
	}

	sc := q.(map[string][]api.ResourceCapacity)

	srvCaps := make([]api.ServiceCapacities, 0, len(sc))
	for srv, resList := range sc {
		resCaps := make([]api.ResourceCapacity, 0, len(resList))
		for _, r := range resList {
			resCaps = append(resCaps, api.ResourceCapacity{
				Name:     r.Name,
				Capacity: r.Capacity,
				Unit:     r.Unit,
				Comment:  r.Comment,
			})
		}

		srvCaps = append(srvCaps, api.ServiceCapacities{
			Type:      srv,
			Resources: resCaps,
		})
	}

	c.IsList = false
	c.Result = clusters.Update(limesClient, c.ID, clusters.UpdateOpts{Services: srvCaps})
	if c.Result.Err != nil {
		kingpin.Fatalf("could not set new capacities for cluster: %v", c.Result.Err)
	}
}

// set updates the resource quota(s) for a domain.
func (d *Domain) set(q interface{}) {
	provider, err := clientconfig.AuthenticatedClient(nil)
	if err != nil {
		kingpin.Fatalf("can not connect to OpenStack: %v", err)
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

	sq := q.(api.ServiceQuotas)

	d.IsList = false
	d.Result = domains.Update(limesClient, d.ID, domains.UpdateOpts{Services: sq})
	if d.Result.Err != nil {
		kingpin.Fatalf("could not set new quota(s) for domain: %v", d.Result.Err)
	}
}

// set updates the resource quota(s) for a project within a specific domain.
func (p *Project) set(q interface{}) {
	provider, err := clientconfig.AuthenticatedClient(nil)
	if err != nil {
		kingpin.Fatalf("can not connect to OpenStack: %v", err)
	}
	identityClient, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		kingpin.Fatalf("could not initialize identity client: %v", err)
	}
	limesClient, err := resources.NewLimesV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		kingpin.Fatalf("could not initialize Limes client: %v", err)
	}

	// if no domain ID is given at the command line then find the project from the
	// user's token scope
	if p.DomainID == "" {
		nameOrID := p.ID
		count, err := p.find(identityClient, "")
		if err != nil {
			kingpin.Fatalf("could not find project in token scope: %v", err)
		}
		if nameOrID == p.Name && count > 1 {
			kingpin.Fatalf("more than one project exists with the name %v", p.Name)
		}

		// if project was found then get the domain name
		if p.DomainID != "" {
			result := idD.Get(identityClient, p.DomainID)
			if result.Err != nil {
				kingpin.Fatalf("could not get project's domain info: %v", result.Err)
			}
			d, err := result.Extract()
			if err != nil {
				kingpin.Fatalf("could not get project's domain info: %v", result.Err)
			}
			p.DomainName = d.Name
		}
	} else {
		nameOrID := p.DomainID
		d := &Domain{ID: p.DomainID}
		count, err := d.find(identityClient)
		if err != nil {
			kingpin.Fatalf("could not find project's domain in token scope: %v", err)
		}
		if nameOrID == d.Name && count > 1 {
			kingpin.Fatalf("more than one domain exists with the name %v", d.Name)
		}

		_, err = p.find(identityClient, d.ID)
		if err != nil {
			kingpin.Fatalf("could not find project in token scope: %v", err)
		}
		p.DomainID = d.ID
		p.DomainName = d.Name
	}

	sq := q.(api.ServiceQuotas)

	p.IsList = false
	p.Result = projects.Update(limesClient, p.DomainID, p.ID, projects.UpdateOpts{Services: sq})
	if p.Result.Err != nil {
		kingpin.Fatalf("could not set new quota(s) for project: %v", p.Result.Err)
	}
}
