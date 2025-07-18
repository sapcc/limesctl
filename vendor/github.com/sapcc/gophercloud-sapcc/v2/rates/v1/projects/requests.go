// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package projects

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/sapcc/go-api-declarations/limes"
)

// ReadOptsBuilder allows extensions to add additional parameters to the Get/List requests.
type ReadOptsBuilder interface {
	ToProjectReadParams() (map[string]string, string, error)
}

// ReadOpts contains parameters for filtering a Get/List request.
type ReadOpts struct {
	Services []limes.ServiceType `q:"service"`
	Areas    []string            `q:"area"`
}

// ToProjectReadParams formats a ReadOpts into a map of headers and a query string.
func (opts ReadOpts) ToProjectReadParams() (headers map[string]string, queryString string, err error) {
	h, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return nil, "", err
	}
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return nil, "", err
	}
	return h, q.String(), nil
}

// List enumerates the projects in a specific domain.
func List(ctx context.Context, c *gophercloud.ServiceClient, domainID string, opts ReadOptsBuilder) (r CommonResult) {
	url := listURL(c, domainID)
	headers := make(map[string]string)
	if opts != nil {
		h, q, err := opts.ToProjectReadParams()
		if err != nil {
			r.Err = err
			return
		}
		headers = h
		url += q
	}

	resp, err := c.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: headers,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves details on a single project, by ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, domainID, projectID string, opts ReadOptsBuilder) (r CommonResult) {
	url := getURL(c, domainID, projectID)
	headers := make(map[string]string)
	if opts != nil {
		h, q, err := opts.ToProjectReadParams()
		if err != nil {
			r.Err = err
			return
		}
		headers = h
		url += q
	}

	resp, err := c.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: headers,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
