// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package projects

import (
	"github.com/gophercloud/gophercloud/v2"
	limesrates "github.com/sapcc/go-api-declarations/limes/rates"
)

// CommonResult is the result of a Get/List operation. Call its appropriate
// Extract method to interpret it as a Project or a slice of Projects.
type CommonResult struct {
	gophercloud.Result
}

// ExtractProjects interprets a CommonResult as a slice of Projects.
func (r CommonResult) ExtractProjects() ([]limesrates.ProjectReport, error) {
	var s struct {
		Projects []limesrates.ProjectReport `json:"projects"`
	}

	err := r.ExtractInto(&s)
	return s.Projects, err
}

// Extract interprets a CommonResult as a Project.
func (r CommonResult) Extract() (*limesrates.ProjectReport, error) {
	var s struct {
		Project *limesrates.ProjectReport `json:"project"`
	}
	err := r.ExtractInto(&s)
	return s.Project, err
}
