// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package clusters

import (
	"github.com/gophercloud/gophercloud/v2"
	limesrates "github.com/sapcc/go-api-declarations/limes/rates"
)

// CommonResult is the result of a Get/List operation. Call its appropriate
// Extract method to interpret it as a Cluster or a slice of Clusters.
type CommonResult struct {
	gophercloud.Result
}

// Extract interprets a CommonResult as a Cluster.
func (r CommonResult) Extract() (*limesrates.ClusterReport, error) {
	var s struct {
		Cluster *limesrates.ClusterReport `json:"cluster"`
	}
	err := r.ExtractInto(&s)
	return s.Cluster, err
}
