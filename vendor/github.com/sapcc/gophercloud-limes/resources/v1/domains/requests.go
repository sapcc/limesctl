// Package domains provides interaction with Limes at the domain hierarchical level.
package domains

import (
	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/limes"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToDomainListParams() (map[string]string, string, error)
}

// ListOpts contains parameters for filtering a List request.
type ListOpts struct {
	Cluster  string `h:"X-Limes-Cluster-Id"`
	Area     string `q:"area"`
	Service  string `q:"service"`
	Resource string `q:"resource"`
}

// ToDomainListParams formats a ListOpts into a map of headers and a query string.
func (opts ListOpts) ToDomainListParams() (map[string]string, string, error) {
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

	_, r.Err = c.Get(url, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: headers,
	})
	return
}

// GetOptsBuilder allows extensions to add additional parameters to the Get request.
type GetOptsBuilder interface {
	ToDomainGetParams() (map[string]string, string, error)
}

// GetOpts contains parameters for filtering a Get request.
type GetOpts struct {
	Cluster  string `h:"X-Limes-Cluster-Id"`
	Area     string `q:"area"`
	Service  string `q:"service"`
	Resource string `q:"resource"`
}

// ToDomainGetParams formats a GetOpts into a map of headers and a query string.
func (opts GetOpts) ToDomainGetParams() (map[string]string, string, error) {
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

	_, r.Err = c.Get(url, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: headers,
	})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToDomainUpdateMap() (map[string]string, map[string]interface{}, error)
}

// UpdateOpts contains parameters to update a domain.
type UpdateOpts struct {
	Cluster  string             `h:"X-Limes-Cluster-Id"`
	Services limes.QuotaRequest `json:"services"`
}

// ToDomainUpdateMap formats a UpdateOpts into a map of headers and a request body.
func (opts UpdateOpts) ToDomainUpdateMap() (map[string]string, map[string]interface{}, error) {
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
func Update(c *gophercloud.ServiceClient, domainID string, opts UpdateOptsBuilder) error {
	url := updateURL(c, domainID)
	h, b, err := opts.ToDomainUpdateMap()
	if err != nil {
		return err
	}
	_, err = c.Put(url, b, nil, &gophercloud.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: h,
	})
	return err
}
