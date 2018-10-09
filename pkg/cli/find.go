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
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/pagination"
)

// find finds a specific domain within the token scope.
// Use if domain name is given instead of a UUID, or if not sure whether the given ID is a name or an UUID.
func (d *Domain) find(client *gophercloud.ServiceClient) (int, error) {
	var nameCount int
	pager := domains.List(client, domains.ListOpts{})
	if pager.Err != nil {
		return 0, pager.Err
	}
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		domainList, err := domains.ExtractDomains(page)
		if err != nil {
			return false, err
		}

		for _, dInList := range domainList {
			if dInList.ID == d.ID || dInList.Name == d.ID {
				d.ID = dInList.ID
				d.Name = dInList.Name
			}
			if dInList.Name == d.Name {
				nameCount++
			}
		}

		return true, nil
	})
	if err != nil {
		return 0, err
	}

	return nameCount, nil
}

// find finds a project within a specific domain in the token scope. If no domain ID is specified then
// it enumerates all the projects within the token scope and finds one in that list.
// Use if project name is given instead of a UUID, or if not sure whether the given ID is a name or an UUID.
func (p *Project) find(client *gophercloud.ServiceClient, findInDomain string) (int, error) {
	var nameCount int
	pager := projects.List(client, projects.ListOpts{})
	if pager.Err != nil {
		return 0, pager.Err
	}
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		projectList, err := projects.ExtractProjects(page)
		if err != nil {
			return false, err
		}

		for _, pInList := range projectList {
			if findInDomain == "" {
				if pInList.ID == p.ID || pInList.Name == p.ID {
					p.ID = pInList.ID
					p.Name = pInList.Name
					p.DomainID = pInList.DomainID
				}
				if pInList.Name == p.Name {
					nameCount++
				}
			} else {
				if (pInList.ID == p.ID || pInList.Name == p.ID) && pInList.DomainID == findInDomain {
					p.ID = pInList.ID
					p.Name = pInList.Name

					return false, nil
				}
			}
		}

		return true, nil
	})
	if err != nil {
		return 0, err
	}

	return nameCount, nil
}
