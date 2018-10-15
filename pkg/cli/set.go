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
	"fmt"

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
	"github.com/sapcc/limes/pkg/limes"
)

// set updates the resource capacities for a cluster.
func (c *Cluster) set(q *Quotas) {
	provider, err := clientconfig.AuthenticatedClient(nil)
	if err != nil {
		kingpin.Fatalf("can not connect to OpenStack: %v", err)
	}
	limesClient, err := resources.NewLimesV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		kingpin.Fatalf("could not initialize Limes client: %v", err)
	}

	srvCaps := make([]api.ServiceCapacities, 0, len(*q))
	for srv, resList := range *q {
		resCaps := make([]api.ResourceCapacity, 0, len(resList))
		for _, r := range resList {
			resCaps = append(resCaps, api.ResourceCapacity{
				Name:     r.Name,
				Capacity: r.Value,
				Unit:     &r.Unit,
				Comment:  r.Comment,
			})
		}

		srvCaps = append(srvCaps, api.ServiceCapacities{
			Type:      srv,
			Resources: resCaps,
		})
	}

	err = clusters.Update(limesClient, c.ID, clusters.UpdateOpts{Services: srvCaps})
	if err != nil {
		kingpin.Fatalf("could not set new capacities for cluster: %v", err)
	}
}

// set updates the resource quota(s) for a domain.
func (d *Domain) set(q *Quotas) {
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

	sq := makeServiceQuotas(q)

	err = domains.Update(limesClient, d.ID, domains.UpdateOpts{Services: sq})
	if err != nil {
		kingpin.Fatalf("could not set new quota(s) for domain: %v", err)
	}
}

// set updates the resource quota(s) for a project within a specific domain.
func (p *Project) set(q *Quotas) {
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

	sq := makeServiceQuotas(q)

	respBody, err := projects.Update(limesClient, p.DomainID, p.ID, projects.UpdateOpts{Services: sq})
	if err != nil {
		kingpin.Fatalf("could not set new quota(s) for project: %v", err)
	}

	if respBody != nil {
		fmt.Println(string(respBody))
	}
}

func makeServiceQuotas(q *Quotas) api.ServiceQuotas {
	sq := make(api.ServiceQuotas)

	for srv, resList := range *q {
		sq[srv] = make(api.ResourceQuotas)

		for _, r := range resList {
			sq[srv][r.Name] = limes.ValueWithUnit{
				Value: uint64(r.Value),
				Unit:  r.Unit,
			}
		}
	}

	return sq
}
