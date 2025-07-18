// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package projects

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/sapcc/go-api-declarations/limes"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToProjectListParams() (map[string]string, string, error)
}

// ListOpts contains parameters for filtering a List request.
type ListOpts struct {
	Detail    bool                          `q:"detail"`
	Areas     []string                      `q:"area"`
	Services  []limes.ServiceType           `q:"service"`
	Resources []limesresources.ResourceName `q:"resource"`
}

// ToProjectListParams formats a ListOpts into a map of headers and a query string.
func (opts ListOpts) ToProjectListParams() (headers map[string]string, queryString string, err error) {
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
func List(ctx context.Context, c *gophercloud.ServiceClient, domainID string, opts ListOptsBuilder) (r CommonResult) {
	url := listURL(c, domainID)
	headers := make(map[string]string)
	if opts != nil {
		h, q, err := opts.ToProjectListParams()
		if err != nil {
			r.Err = err
			return
		}
		headers = h
		url += q
	}

	resp, err := c.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{ //nolint:bodyclose // already closed by gophercloud
		MoreHeaders: headers,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetOptsBuilder allows extensions to add additional parameters to the Get request.
type GetOptsBuilder interface {
	ToProjectGetParams() (map[string]string, string, error)
}

// GetOpts contains parameters for filtering a Get request.
type GetOpts struct {
	Detail    bool                          `q:"detail"`
	Areas     []string                      `q:"area"`
	Services  []limes.ServiceType           `q:"service"`
	Resources []limesresources.ResourceName `q:"resource"`
}

// ToProjectGetParams formats a GetOpts into a map of headers and a query string.
func (opts GetOpts) ToProjectGetParams() (headers map[string]string, queryString string, err error) {
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

// Get retrieves details on a single project, by ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, domainID, projectID string, opts GetOptsBuilder) (r CommonResult) {
	url := getURL(c, domainID, projectID)
	headers := make(map[string]string)
	if opts != nil {
		h, q, err := opts.ToProjectGetParams()
		if err != nil {
			r.Err = err
			return
		}
		headers = h
		url += q
	}

	resp, err := c.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{ //nolint:bodyclose // already closed by gophercloud
		MoreHeaders: headers,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Sync schedules a sync task that pulls a project's data from the backing services
// into Limes' local database.
func Sync(ctx context.Context, c *gophercloud.ServiceClient, domainID, projectID string) (r SyncResult) {
	url := syncURL(c, domainID, projectID)
	resp, err := c.Post(ctx, url, nil, nil, &gophercloud.RequestOpts{ //nolint:bodyclose // already closed by gophercloud
		OkCodes: []int{http.StatusAccepted},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
