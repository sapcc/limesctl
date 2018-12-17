package domains

import (
	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/limes"
)

// CommonResult is the result of a Get/List operation. Call its appropriate
// Extract method to interpret it as a Domain or a slice of Domains.
type CommonResult struct {
	gophercloud.Result
}

// ExtractDomains interprets a CommonResult as a slice of Domains.
func (r CommonResult) ExtractDomains() ([]limes.DomainReport, error) {
	var s struct {
		Domains []limes.DomainReport `json:"domains"`
	}

	err := r.ExtractInto(&s)
	return s.Domains, err
}

// Extract interprets a CommonResult as a Domain.
func (r CommonResult) Extract() (*limes.DomainReport, error) {
	var s struct {
		Domain *limes.DomainReport `json:"domain"`
	}
	err := r.ExtractInto(&s)
	return s.Domain, err
}
