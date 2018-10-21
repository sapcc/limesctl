package projects

import (
	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/limes/pkg/reports"
)

// CommonResult is the result of a Get/List operation. Call its appropriate
// Extract method to interpret it as a Project or a slice of Projects.
type CommonResult struct {
	gophercloud.Result
}

// ExtractProjects interprets a CommonResult as a slice of Projects.
func (r CommonResult) ExtractProjects() ([]reports.Project, error) {
	var s struct {
		Projects []reports.Project `json:"projects"`
	}

	err := r.ExtractInto(&s)
	return s.Projects, err
}

// Extract interprets a CommonResult as a Project.
func (r CommonResult) Extract() (*reports.Project, error) {
	var s struct {
		Project *reports.Project `json:"project"`
	}
	err := r.ExtractInto(&s)
	return s.Project, err
}
