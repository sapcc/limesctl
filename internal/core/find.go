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
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud"
	gopherdomains "github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	gopherprojects "github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/sapcc/gophercloud-limes/resources/v1/domains"
	"github.com/sapcc/gophercloud-limes/resources/v1/projects"
	"github.com/sapcc/limesctl/internal/auth"
	"github.com/sapcc/limesctl/internal/errors"
)

// FindDomain uses the user's input (name/UUID) to find a specific domain
// within the token scope.
// Different strategies are tried in a chronological order to find the relevant
// domain in the most efficient way possible.
func FindDomain(identityV3, limesV1 *gophercloud.ServiceClient, userInput, clusterID string) (*Domain, error) {
	// Strategy 1: if clusterID is given then userInput is assumed to be an ID
	if clusterID != "" {
		return findDomainInCluster(limesV1, userInput, clusterID)
	}

	// Strategy 2: check if the domain is mentioned in the current token scope
	token, err := auth.CurrentToken(identityV3)
	if err == nil {
		d1 := token.Domain
		if d1.ID != "" && (d1.Name == userInput || d1.ID == userInput) {
			return &Domain{Name: d1.Name, ID: d1.ID}, nil
		}
		d2 := token.Project.Domain
		if d2.ID != "" && (d2.Name == userInput || d2.ID == userInput) {
			return &Domain{Name: d2.Name, ID: d2.ID}, nil
		}
	}

	// Strategy 3: assume that userInput is an ID
	d, err := gopherdomains.Get(identityV3, userInput).Extract()
	if err == nil {
		return &Domain{Name: d.Name, ID: d.ID}, nil
	}

	// Strategy 4: assume userInput is a name and do a domain listing
	// restricted to that specific name
	page, err := gopherdomains.List(identityV3, gopherdomains.ListOpts{Name: userInput}).AllPages()
	if err != nil {
		return nil, err
	}
	dList, err := gopherdomains.ExtractDomains(page)
	if err != nil {
		return nil, err
	}
	if len(dList) > 1 {
		return nil, fmt.Errorf("more than one domain exists with the name %q", userInput)
	}
	if len(dList) != 0 && dList[0].ID != "" {
		return &Domain{Name: dList[0].Name, ID: dList[0].ID}, nil
	}

	// at this point all strategies have failed
	return nil, errors.New("domain not found")
}

func findDomainInCluster(limesV1 *gophercloud.ServiceClient, domainID, clusterID string) (*Domain, error) {
	var s struct {
		Domain struct {
			UUID string `json:"id"`
			Name string `json:"name"`
		} `json:"domain"`
	}

	// gophercloud doesn't support projects across different clusters therefore
	// gophercloud-limes is used here
	err := domains.Get(limesV1, domainID, domains.GetOpts{Cluster: clusterID}).ExtractInto(&s)
	if err != nil {
		return nil, fmt.Errorf("could not find domain: %v", err)
	}

	return &Domain{ID: s.Domain.UUID, Name: s.Domain.Name}, nil
}

// FindProject uses the user's input (name/UUID) to find a specific project
// within the token scope.
// Different strategies are tried in a chronological order to find the relevant
// project in the most efficient way possible.
func FindProject(identityV3, limesV1 *gophercloud.ServiceClient, userInputProject, userInputDomain, clusterID string) (*Project, error) {
	// Strategy 1: if clusterID is given then userInputs are assumed to be IDs
	if clusterID != "" {
		return findProjectInCluster(limesV1, userInputProject, userInputDomain, clusterID)
	}

	// Strategy 2: check if the project is mentioned in the current token scope
	token, err := auth.CurrentToken(identityV3)
	if err == nil {
		p := token.Project
		if p.ID != "" && (p.ID == userInputProject || p.Name == userInputProject) {
			d1 := token.Domain
			if d1.ID != "" && (d1.Name == userInputDomain || d1.ID == userInputDomain) {
				return &Project{
					ID:         p.ID,
					Name:       p.Name,
					DomainID:   d1.ID,
					DomainName: d1.Name,
				}, nil
			}
			d2 := token.Project.Domain
			if d2.ID != "" && (d2.Name == userInputDomain || d2.ID == userInputDomain) {
				return &Project{
					ID:         p.ID,
					Name:       p.Name,
					DomainID:   d2.ID,
					DomainName: d2.Name,
				}, nil
			}
		}
	}

	// Strategy 3: assume that userInputProject is an ID
	p, err := gopherprojects.Get(identityV3, userInputProject).Extract()
	if err == nil {
		// get domain name
		d, err := gopherdomains.Get(identityV3, p.DomainID).Extract()
		if err != nil {
			return nil, fmt.Errorf("could not find project: %v", err)
		}

		return &Project{
			ID:         p.ID,
			Name:       p.Name,
			DomainID:   d.ID,
			DomainName: d.Name,
		}, nil
	}

	// Strategy 4: assume userInputProject is a name and do a project listing
	// restricted to that specific name
	var page pagination.Page
	if userInputDomain != "" {
		d, err := FindDomain(identityV3, limesV1, userInputDomain, clusterID)
		if err == nil {
			page, err = gopherprojects.List(identityV3, gopherprojects.ListOpts{
				Name:     userInputProject,
				DomainID: d.ID,
			}).AllPages()
		}
	} else {
		page, err = gopherprojects.List(identityV3, gopherprojects.ListOpts{Name: userInputProject}).AllPages()
	}
	if err != nil {
		return nil, err
	}

	// this is needed in case the user did not gave a domain ID at input
	// which means we still don't have the domain name
	if p.DomainName == "" {
		d, err := domains.Get(identityV3, p.DomainID).Extract()
		if err == nil {
			p.DomainName = d.Name
		} else if strings.Contains(err.Error(), "Forbidden") {
			//if the user can access the project, but does not have permissions for
			//`openstack domain show`, continue with a bogus domain name (this issue
			//would otherwise completely break limesctl for that user even though
			//they have permissions for the Limes API)
			p.DomainName = "domain-" + p.DomainID
		} else {
			//unexpected error
			return nil, err
		}
	}

	// at this point all strategies have failed
	return nil, errors.New("project not found")
}

func findProjectInCluster(limesV1 *gophercloud.ServiceClient, projectID, domainID, clusterID string) (*Project, error) {
	var s struct {
		Project struct {
			UUID string `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
	}

	// gophercloud doesn't support projects across different clusters therefore
	// gophercloud-limes is used here
	err := projects.Get(limesV1, domainID, projectID, projects.GetOpts{
		Cluster: clusterID}).ExtractInto(&s)
	if err != nil {
		return nil, fmt.Errorf("could not find project: %v", err)
	}

	// get domain name
	d, err := findDomainInCluster(limesV1, domainID, clusterID)
	if err != nil {
		return nil, err
	}

	return &Project{
		ID:         s.Project.UUID,
		Name:       s.Project.Name,
		DomainID:   d.ID,
		DomainName: d.Name,
	}, nil
}
