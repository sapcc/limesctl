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
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/gophercloud-sapcc/resources/v1/projects"
	"github.com/sapcc/limesctl/internal/errors"
)

// set updates the resource quota(s) for a project within a specific domain.
func (p *Project) set(limesV1 *gophercloud.ServiceClient, q Quotas) {
	sq := makeServiceQuotas(q)

	warn, err := projects.Update(limesV1, p.DomainID, p.ID, projects.UpdateOpts{
		Cluster:  p.Filter.Cluster,
		Services: sq,
	}).Extract()
	errors.Handle(err, "could not set new quota(s) for project")

	if warn != nil {
		fmt.Println(string(warn))
	}
}
