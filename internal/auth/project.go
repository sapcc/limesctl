// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	identityprojects "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"

	"github.com/sapcc/limesctl/v3/internal/util"
)

const msgProjectNotFound = "project not found"

// ProjectInfo identifies a specific project.
type ProjectInfo struct {
	ID         string
	DomainID   string
	DomainName string
}

// FindProject tries to find a project using the provided name/ID(s).
func FindProject(ctx context.Context, identityClient *gophercloud.ServiceClient, domainNameOrID, projectNameOrID string) (*ProjectInfo, error) {
	// Strategy 1: find project in current token.
	pInfo := findProjectInCurrentToken(identityClient, domainNameOrID, projectNameOrID)
	if pInfo != nil {
		return pInfo, nil
	} else if projectNameOrID == "" {
		// If no projectNameOrID is provided and we can't find the project from
		// token scope then further strategies are futile.
		return nil, errors.New(msgProjectNotFound)
	}

	// Strategy 2: assume that projectNameOrID is an ID and try to find in
	// Keystone.
	p, err := identityprojects.Get(ctx, identityClient, projectNameOrID).Extract()
	if err == nil {
		if p.IsDomain {
			return nil, errors.New("the given ID belongs to a domain, usage instructions: limectl domain --help")
		}
		if p.ID != "" {
			var dName string
			dName, err = FindDomainName(ctx, identityClient, p.DomainID)
			if err != nil {
				return nil, util.WrapError(err, msgProjectNotFound)
			}
			return &ProjectInfo{
				ID:         p.ID,
				DomainID:   p.DomainID,
				DomainName: dName,
			}, nil
		}
	}

	// Strategy 3: at this point we know that projectNameOrID is a name so we
	// do a Keystone project listing and try to find the project.
	opts := identityprojects.ListOpts{Name: projectNameOrID}
	if domainNameOrID != "" {
		var dID string
		dID, err = FindDomainID(ctx, identityClient, domainNameOrID)
		if err != nil {
			return nil, util.WrapError(err, msgProjectNotFound)
		}
		opts.DomainID = dID
	}
	var pList []identityprojects.Project
	page, err := identityprojects.List(identityClient, opts).AllPages(ctx)
	if err == nil {
		pList, err = identityprojects.ExtractProjects(page)
	}
	if err != nil {
		return nil, util.WrapError(err, msgProjectNotFound)
	}
	l := len(pList)
	if l > 1 {
		return nil, fmt.Errorf("more than one project exists with the name %q", projectNameOrID)
	}
	if l == 1 {
		if p := pList[0]; p.ID != "" {
			dName, err := FindDomainName(ctx, identityClient, p.DomainID)
			if err != nil {
				return nil, util.WrapError(err, msgProjectNotFound)
			}
			return &ProjectInfo{
				ID:         p.ID,
				DomainID:   p.DomainID,
				DomainName: dName,
			}, nil
		}
	}

	// All strategies have failed :(
	return nil, errors.New(msgProjectNotFound)
}

func findProjectInCurrentToken(identityClient *gophercloud.ServiceClient, domainNameOrID, projectNameOrID string) *ProjectInfo {
	t, err := currentToken(identityClient)
	if err != nil {
		return nil
	}

	p := t.Project
	d1 := t.Project.Domain
	d2 := t.Domain
	if p.ID == "" {
		return nil
	}

	switch projectNameOrID {
	case "":
		// If no projectNameOrID is provided then we return the info from token.
		if d1.ID != "" {
			return &ProjectInfo{
				ID:         p.ID,
				DomainID:   d1.ID,
				DomainName: d1.Name,
			}
		}
		if d2.ID != "" {
			return &ProjectInfo{
				ID:         p.ID,
				DomainID:   d2.ID,
				DomainName: d2.Name,
			}
		}
	case p.ID, p.Name:
		// Check if token has the given name/ID(s).
		if d1.ID != "" && (domainNameOrID == d1.ID || domainNameOrID == d1.Name) {
			return &ProjectInfo{
				ID:         p.ID,
				DomainID:   d1.ID,
				DomainName: d1.Name,
			}
		}
		if d2.ID != "" && (domainNameOrID == d2.ID || domainNameOrID == d2.Name) {
			return &ProjectInfo{
				ID:         p.ID,
				DomainID:   d2.ID,
				DomainName: d2.Name,
			}
		}
	}

	return nil
}
