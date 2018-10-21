package clusters

import (
	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/limes/pkg/reports"
)

// CommonResult is the result of a Get/List operation. Call its appropriate
// Extract method to interpret it as a Cluster or a slice of Clusters.
type CommonResult struct {
	gophercloud.Result
}

// ExtractClusters interprets a CommonResult as a slice of Clusters.
func (r CommonResult) ExtractClusters() ([]reports.Cluster, error) {
	var s struct {
		Clusters []reports.Cluster `json:"clusters"`
	}

	err := r.ExtractInto(&s)
	return s.Clusters, err
}

// Extract interprets a CommonResult as a Cluster.
func (r CommonResult) Extract() (*reports.Cluster, error) {
	var s struct {
		Cluster *reports.Cluster `json:"cluster"`
	}
	err := r.ExtractInto(&s)
	return s.Cluster, err
}
