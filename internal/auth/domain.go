// Copyright 2020 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	identitydomains "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"

	"github.com/sapcc/limesctl/v3/internal/util"
)

const msgDomainNotFound = "domain not found"

// FindDomainName returns the name of a domain.
func FindDomainName(ctx context.Context, identityClient *gophercloud.ServiceClient, id string) (string, error) {
	d, err := identitydomains.Get(ctx, identityClient, id).Extract()
	switch {
	case err == nil:
		return d.Name, nil
	case strings.Contains(err.Error(), "Forbidden"):
		// If the user can access the project but does not have permissions for
		// `openstack domain show` then return a pseudoname.
		// This issue would otherwise completely break limesctl for that user
		// even though they have permissions for the Limes API.
		n := "domain-" + id
		return n, nil
	default:
		return "", util.WrapError(err, msgDomainNotFound)
	}
}

// FindDomainID tries to find a domain id using the provided nameOrID.
func FindDomainID(ctx context.Context, identityClient *gophercloud.ServiceClient, nameOrID string) (string, error) {
	// Strategy 1: get domain id from current token scope.
	id := getDomainIDFromCurrentToken(identityClient, nameOrID)
	if id != "" {
		return id, nil
	} else if nameOrID == "" {
		// If no nameOrID is provided and we can't find the id from token scope
		// then further strategies are futile.
		return "", errors.New(msgDomainNotFound)
	}

	// Strategy 2: assume that nameOrID is an ID and try to find in Keystone.
	d, err := identitydomains.Get(ctx, identityClient, nameOrID).Extract()
	if err == nil && d.ID != "" {
		return d.ID, nil
	}

	// Strategy 3: at this point we know that nameOrID is a name so we do a
	// Keystone domain listing and try to find the domain.
	var dList []identitydomains.Domain
	page, err := identitydomains.List(identityClient, identitydomains.ListOpts{Name: nameOrID}).AllPages(ctx)
	if err == nil {
		dList, err = identitydomains.ExtractDomains(page)
	}
	if err != nil {
		return "", util.WrapError(err, msgDomainNotFound)
	}
	l := len(dList)
	if l > 1 {
		return "", fmt.Errorf("more than one domain exists with the name %q", nameOrID)
	}
	if l == 1 {
		if id := dList[0].ID; id != "" {
			return id, nil
		}
	}

	// All strategies have failed :(
	return "", errors.New(msgDomainNotFound)
}

func getDomainIDFromCurrentToken(identityClient *gophercloud.ServiceClient, nameOrID string) string {
	t, err := currentToken(identityClient)
	if err != nil {
		return ""
	}

	d1 := t.Domain
	d2 := t.Project.Domain
	if nameOrID == "" {
		// If no nameOrID is provided then we return the id from token.
		if d1.ID != "" {
			return d1.ID
		}
		if d2.ID != "" {
			return d2.ID
		}
	} else {
		// Check if token has the nameOrID.
		if d1.ID != "" && (nameOrID == d1.ID || nameOrID == d1.Name) {
			return d1.ID
		}
		if d2.ID != "" && (nameOrID == d2.ID || nameOrID == d2.Name) {
			return d2.ID
		}
	}
	return ""
}
