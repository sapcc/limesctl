// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

// Package clusters provides interaction with Limes at the cluster hierarchical level.
package clusters

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/sapcc/go-api-declarations/limes"
)

// GetOptsBuilder allows extensions to add additional parameters to the Get request.
type GetOptsBuilder interface {
	ToClusterGetQuery() (string, error)
}

// GetOpts contains parameters for filtering a Get request.
type GetOpts struct {
	Services []limes.ServiceType `q:"service"`
	Areas    []string            `q:"area"`
}

// ToClusterGetQuery formats a GetOpts into a query string.
func (opts GetOpts) ToClusterGetQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get retrieves details on a single cluster, by ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, opts GetOptsBuilder) (r CommonResult) {
	url := getURL(c)
	if opts != nil {
		query, err := opts.ToClusterGetQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := c.Get(ctx, url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
