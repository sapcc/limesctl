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
	"github.com/sapcc/limesctl/pkg/errors"
)

// list retrieves information about all the clusters within the token scope.
func (c *Cluster) list() {
	_, limesV1 := getServiceClients()

	c.IsList = true
	c.Result = clusters.List(limesV1, clusters.ListOpts{
		Area:     c.Filter.Area,
		Service:  c.Filter.Service,
		Resource: c.Filter.Resource,
	})
	errors.Handle(c.Result.Err, "could not list clusters")
}

// list retrieves information about all the domains within the token scope.
func (d *Domain) list() {
	_, limesV1 := getServiceClients()

	d.IsList = true
	d.Result = domains.List(limesV1, domains.ListOpts{
		Cluster:  d.Filter.Cluster,
		Area:     d.Filter.Area,
		Service:  d.Filter.Service,
		Resource: d.Filter.Resource,
	})
	errors.Handle(d.Result.Err, "could not list domains")
}

// list retrieves information about all the projects within a specific domain.
func (p *Project) list() {
	_, limesV1 := getServiceClients()

	p.IsList = true
	p.Result = projects.List(limesV1, p.DomainID, projects.ListOpts{
		Cluster:  p.Filter.Cluster,
		Area:     p.Filter.Area,
		Service:  p.Filter.Service,
		Resource: p.Filter.Resource,
	})
	errors.Handle(p.Result.Err, "could not list projects")
}
