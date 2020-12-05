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

package core

import (
	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/gophercloud-sapcc/resources/v1/projects"
	"github.com/sapcc/limesctl/internal/errors"
)

// list retrieves information about all the projects within a specific domain.
func (p *Project) list(limesV1 *gophercloud.ServiceClient) {
	p.IsList = true
	p.Result = projects.List(limesV1, p.DomainID, projects.ListOpts{
		Cluster:  p.Filter.Cluster,
		Area:     p.Filter.Area,
		Service:  p.Filter.Service,
		Resource: p.Filter.Resource,
	})
	errors.Handle(p.Result.Err, "could not list projects")
}
