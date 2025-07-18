// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package clusters

import (
	"github.com/gophercloud/gophercloud/v2"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
)

// CommonResult is the result of a Get operation. Call its appropriate
// Extract method to interpret it as a Cluster.
type CommonResult struct {
	gophercloud.Result
}

// Extract interprets a CommonResult as a Cluster.
func (r CommonResult) Extract() (*limesresources.ClusterReport, error) {
	var s struct {
		Cluster *limesresources.ClusterReport `json:"cluster"`
	}
	err := r.ExtractInto(&s)
	return s.Cluster, err
}
