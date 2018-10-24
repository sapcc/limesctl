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
	"github.com/sapcc/gophercloud-limes/resources/v1/clusters"
	"github.com/sapcc/gophercloud-limes/resources/v1/domains"
	"github.com/sapcc/gophercloud-limes/resources/v1/projects"
)

// get retrieves information about a single cluster.
func (c *Cluster) get() {
	_, limesV1 := getServiceClients()

	c.Result = clusters.Get(limesV1, c.ID, clusters.GetOpts{
		Area:     c.Opts.Area,
		Service:  c.Opts.Service,
		Resource: c.Opts.Resource,
	})
	handleError("could not get cluster", c.Result.Err)
}

// get retrieves information about a single domain.
func (d *Domain) get() {
	_, limesV1 := getServiceClients()

	d.Result = domains.Get(limesV1, d.ID, domains.GetOpts{
		Cluster:  d.Opts.Cluster,
		Area:     d.Opts.Area,
		Service:  d.Opts.Service,
		Resource: d.Opts.Resource,
	})
	handleError("could not get domain", d.Result.Err)
}

// get retrieves information about a single project within a specific domain.
func (p *Project) get() {
	_, limesV1 := getServiceClients()

	p.Result = projects.Get(limesV1, p.DomainID, p.ID, projects.GetOpts{
		Cluster:  p.Opts.Cluster,
		Area:     p.Opts.Area,
		Service:  p.Opts.Service,
		Resource: p.Opts.Resource,
	})
	handleError("could not get project", p.Result.Err)
}
