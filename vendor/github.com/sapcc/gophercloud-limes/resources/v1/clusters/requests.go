// Package clusters provides interaction with Limes at the cluster hierarchical level.
package clusters

import (
	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/limes/pkg/api"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToClusterListQuery() (string, error)
}

// ListOpts contains parameters for filtering a List request
type ListOpts struct {
	Detail   bool   `q:"detail"`
	Local    bool   `q:"local"`
	Area     string `q:"area"`
	Service  string `q:"service"`
	Resource string `q:"resource"`
}

// ToClusterListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the clusters to which the current token has access.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) (r CommonResult) {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToClusterListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// GetOptsBuilder allows extensions to add additional parameters to the Get request.
type GetOptsBuilder interface {
	ToClusterGetQuery() (string, error)
}

// GetOpts contains parameters for filtering a Get request.
type GetOpts struct {
	Detail   bool   `q:"detail"`
	Local    bool   `q:"local"`
	Area     string `q:"area"`
	Service  string `q:"service"`
	Resource string `q:"resource"`
}

// ToClusterGetQuery formats a GetOpts into a query string.
func (opts GetOpts) ToClusterGetQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get retrieves details on a single cluster, by ID.
func Get(c *gophercloud.ServiceClient, clusterID string, opts GetOptsBuilder) (r CommonResult) {
	url := getURL(c, clusterID)
	if opts != nil {
		query, err := opts.ToClusterGetQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToClusterUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains parameters to update a cluster.
type UpdateOpts struct {
	Services []api.ServiceCapacities `json:"services"`
}

// ToClusterUpdateMap formats a UpdateOpts into a request body.
func (opts UpdateOpts) ToClusterUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "cluster")
}

// Update modifies the attributes of a cluster.
func Update(c *gophercloud.ServiceClient, clusterID string, opts UpdateOptsBuilder) error {
	url := updateURL(c, clusterID)
	b, err := opts.ToClusterUpdateMap()
	if err != nil {
		return err
	}
	_, err = c.Put(url, b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return err
}
