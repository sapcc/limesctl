package auth

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	identityprojects "github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/pkg/errors"
)

const msgProjectNotFound = "project not found"

// ProjectInfo identifies a specific project.
type ProjectInfo struct {
	ID         string
	DomainID   string
	DomainName string
}

// FindProject tries to find a project using the provided name/ID(s).
func FindProject(identityClient *gophercloud.ServiceClient, domainNameOrID, projectNameOrID string) (*ProjectInfo, error) {
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
	p, err := identityprojects.Get(identityClient, projectNameOrID).Extract()
	if err == nil && p.ID != "" {
		dName, err := FindDomainName(identityClient, p.DomainID)
		if err != nil {
			return nil, errors.Wrap(err, msgProjectNotFound)
		}
		return &ProjectInfo{
			ID:         p.ID,
			DomainID:   p.DomainID,
			DomainName: dName,
		}, nil
	}

	// Strategy 3: at this point we know that projectNameOrID is a name so we
	// do a Keystone project listing and try to find the project.
	opts := identityprojects.ListOpts{Name: projectNameOrID}
	if domainNameOrID != "" {
		dID, err := FindDomainID(identityClient, domainNameOrID)
		if err != nil {
			return nil, errors.Wrap(err, msgProjectNotFound)
		}
		opts.DomainID = dID
	}
	var pList []identityprojects.Project
	page, err := identityprojects.List(identityClient, opts).AllPages()
	if err == nil {
		pList, err = identityprojects.ExtractProjects(page)
	}
	if err != nil {
		return nil, errors.Wrap(err, msgProjectNotFound)
	}
	l := len(pList)
	if l > 1 {
		return nil, fmt.Errorf("more than one project exists with the name %q", projectNameOrID)
	}
	if l == 1 {
		if p := pList[0]; p.ID != "" {
			dName, err := FindDomainName(identityClient, p.DomainID)
			if err != nil {
				return nil, errors.Wrap(err, msgProjectNotFound)
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

	if projectNameOrID == "" {
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
	} else if projectNameOrID == p.ID || projectNameOrID == p.Name {
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
