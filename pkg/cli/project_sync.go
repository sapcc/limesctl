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
	"github.com/sapcc/gophercloud-limes/resources/v1/projects"
)

// RunSyncTask schedules a sync job that pulls quota and usage data for a project from
// the backing services into Limes' local database.
func RunSyncTask(p *Project) {
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
	}

	err = projects.Sync(limesClient, p.DomainID, p.ID)
	if err != nil {
		kingpin.Fatalf("could not sync project: %v", err)
	}
}
