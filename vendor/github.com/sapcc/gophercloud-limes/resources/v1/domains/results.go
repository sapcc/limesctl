package domains

import (
	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/limes/pkg/reports"
)

// CommonResult is the result of a Get/List operation. Call its appropriate
// Extract method to interpret it as a Domain or a slice of Domains.
type CommonResult struct {
	gophercloud.Result
}

// ExtractDomains interprets a CommonResult as a slice of Domains.
func (r CommonResult) ExtractDomains() ([]reports.Domain, error) {
	var s struct {
		Domains []reports.Domain `json:"domains"`
	}

	err := r.ExtractInto(&s)
	return s.Domains, err
}

// Extract interprets a CommonResult as a Domain.
func (r CommonResult) Extract() (*reports.Domain, error) {
	var s struct {
		Domain *reports.Domain `json:"domain"`
	}
	err := r.ExtractInto(&s)
	return s.Domain, err
}
