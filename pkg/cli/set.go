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

	"github.com/sapcc/gophercloud-limes/resources/v1/clusters"
	"github.com/sapcc/gophercloud-limes/resources/v1/domains"
	"github.com/sapcc/gophercloud-limes/resources/v1/projects"
)

// set updates the resource capacities for a cluster.
func (c *Cluster) set(q *Quotas) {
	_, limesV1 := getServiceClients()

	sc := makeServiceCapacities(q)

	err := clusters.Update(limesV1, c.ID, clusters.UpdateOpts{Services: sc})
	handleError("could not set new capacities for cluster", err)
}

// set updates the resource quota(s) for a domain.
func (d *Domain) set(q *Quotas) {
	_, limesV1 := getServiceClients()

	sq := makeServiceQuotas(q)

	err := domains.Update(limesV1, d.ID, domains.UpdateOpts{
		Cluster:  d.Opts.Cluster,
		Services: sq,
	})
	handleError("could not set new quota(s) for domain", err)
}

// set updates the resource quota(s) for a project within a specific domain.
func (p *Project) set(q *Quotas) {
	_, limesV1 := getServiceClients()

	sq := makeServiceQuotas(q)

	respBody, err := projects.Update(limesV1, p.DomainID, p.ID, projects.UpdateOpts{
		Cluster:  p.Opts.Cluster,
		Services: sq,
	})
	handleError("could not set new quota(s) for project", err)

	if respBody != nil {
		fmt.Println(string(respBody))
	}
}
