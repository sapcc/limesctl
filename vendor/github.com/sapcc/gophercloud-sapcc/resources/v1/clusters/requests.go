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

// Package clusters provides interaction with Limes at the cluster hierarchical level.
package clusters

import (
	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/go-api-declarations/limes"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
)

// GetOptsBuilder allows extensions to add additional parameters to the Get request.
type GetOptsBuilder interface {
	ToClusterGetQuery() (string, error)
}

// GetOpts contains parameters for filtering a Get request.
type GetOpts struct {
	Detail    bool                          `q:"detail"`
	Areas     []string                      `q:"area"`
	Services  []limes.ServiceType           `q:"service"`
	Resources []limesresources.ResourceName `q:"resource"`
}

// ToClusterGetQuery formats a GetOpts into a query string.
func (opts GetOpts) ToClusterGetQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get retrieves details on a single cluster, by ID.
func Get(c *gophercloud.ServiceClient, opts GetOptsBuilder) (r CommonResult) {
	url := getURL(c)
	if opts != nil {
		query, err := opts.ToClusterGetQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := c.Get(url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
