package clusters

import (
	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/limes"
)

// CommonResult is the result of a Get/List operation. Call its appropriate
// Extract method to interpret it as a Cluster or a slice of Clusters.
type CommonResult struct {
	gophercloud.Result
}

// ExtractClusters interprets a CommonResult as a slice of Clusters.
func (r CommonResult) ExtractClusters() ([]limes.ClusterReport, error) {
	var s struct {
		Clusters []limes.ClusterReport `json:"clusters"`
	}

	err := r.ExtractInto(&s)
	return s.Clusters, err
}

// Extract interprets a CommonResult as a Cluster.
func (r CommonResult) Extract() (*limes.ClusterReport, error) {
	var s struct {
		Cluster *limes.ClusterReport `json:"cluster"`
	}
	err := r.ExtractInto(&s)
	return s.Cluster, err
}
