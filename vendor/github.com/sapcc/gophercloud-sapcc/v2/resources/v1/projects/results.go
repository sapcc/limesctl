// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package projects

import (
	"github.com/gophercloud/gophercloud/v2"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
)

// CommonResult is the result of a Get/List operation. Call its appropriate
// Extract method to interpret it as a Project or a slice of Projects.
type CommonResult struct {
	gophercloud.Result
}

// SyncResult is the result of an Sync operation. Call its appropriate
// ExtractErr method to extract the error from the result.
type SyncResult struct {
	gophercloud.ErrResult
}

// ExtractProjects interprets a CommonResult as a slice of Projects.
func (r CommonResult) ExtractProjects() ([]limesresources.ProjectReport, error) {
	var s struct {
		Projects []limesresources.ProjectReport `json:"projects"`
	}

	err := r.ExtractInto(&s)
	return s.Projects, err
}

// Extract interprets a CommonResult as a Project.
func (r CommonResult) Extract() (*limesresources.ProjectReport, error) {
	var s struct {
		Project *limesresources.ProjectReport `json:"project"`
	}
	err := r.ExtractInto(&s)
	return s.Project, err
}
