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

// Package domains provides interaction with Limes at the domain hierarchical level.
package domains

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/go-api-declarations/limes"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToDomainListParams() (map[string]string, string, error)
}

// ListOpts contains parameters for filtering a List request.
type ListOpts struct {
	Areas     []string                      `q:"area"`
	Services  []limes.ServiceType           `q:"service"`
	Resources []limesresources.ResourceName `q:"resource"`
}

// ToDomainListParams formats a ListOpts into a map of headers and a query string.
func (opts ListOpts) ToDomainListParams() (headers map[string]string, queryString string, err error) {
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

// List enumerates the domains to which the current token has access.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) (r CommonResult) {
	url := listURL(c)
	headers := make(map[string]string)
	if opts != nil {
		h, q, err := opts.ToDomainListParams()
		if err != nil {
			r.Err = err
			return
		}
		headers = h
		url += q
	}

	resp, err := c.Get(url, &r.Body, &gophercloud.RequestOpts{ //nolint:bodyclose // already closed by gophercloud
		MoreHeaders: headers,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetOptsBuilder allows extensions to add additional parameters to the Get request.
type GetOptsBuilder interface {
	ToDomainGetParams() (map[string]string, string, error)
}

// GetOpts contains parameters for filtering a Get request.
type GetOpts struct {
	Areas     []string                      `q:"area"`
	Services  []limes.ServiceType           `q:"service"`
	Resources []limesresources.ResourceName `q:"resource"`
}

// ToDomainGetParams formats a GetOpts into a map of headers and a query string.
func (opts GetOpts) ToDomainGetParams() (headers map[string]string, quersString string, err error) {
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

// Get retrieves details on a single domain, by ID.
func Get(c *gophercloud.ServiceClient, domainID string, opts GetOptsBuilder) (r CommonResult) {
	url := getURL(c, domainID)
	headers := make(map[string]string)
	if opts != nil {
		h, q, err := opts.ToDomainGetParams()
		if err != nil {
			r.Err = err
			return
		}
		headers = h
		url += q
	}

	resp, err := c.Get(url, &r.Body, &gophercloud.RequestOpts{ //nolint:bodyclose // already closed by gophercloud
		MoreHeaders: headers,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToDomainUpdateMap() (map[string]string, map[string]interface{}, error)
}

// UpdateOpts contains parameters to update a domain.
type UpdateOpts struct {
	Services limesresources.QuotaRequest `json:"services"`
}

// ToDomainUpdateMap formats a UpdateOpts into a map of headers and a request body.
func (opts UpdateOpts) ToDomainUpdateMap() (headers map[string]string, requestBody map[string]interface{}, err error) {
	h, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return nil, nil, err
	}

	b, err := gophercloud.BuildRequestBody(opts, "domain")
	if err != nil {
		return nil, nil, err
	}

	return h, b, nil
}

// Update modifies the attributes of a domain.
func Update(c *gophercloud.ServiceClient, domainID string, opts UpdateOptsBuilder) (r UpdateResult) {
	url := updateURL(c, domainID)
	h, b, err := opts.ToDomainUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(url, b, nil, &gophercloud.RequestOpts{ //nolint:bodyclose // already closed by gophercloud
		OkCodes:     []int{http.StatusAccepted},
		MoreHeaders: h,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
