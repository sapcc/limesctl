// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package domains

import (
	"github.com/gophercloud/gophercloud/v2"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
)

// CommonResult is the result of a Get/List operation. Call its appropriate
// Extract method to interpret it as a Domain or a slice of Domains.
type CommonResult struct {
	gophercloud.Result
}

// ExtractDomains interprets a CommonResult as a slice of Domains.
func (r CommonResult) ExtractDomains() ([]limesresources.DomainReport, error) {
	var s struct {
		Domains []limesresources.DomainReport `json:"domains"`
	}

	err := r.ExtractInto(&s)
	return s.Domains, err
}

// Extract interprets a CommonResult as a Domain.
func (r CommonResult) Extract() (*limesresources.DomainReport, error) {
	var s struct {
		Domain *limesresources.DomainReport `json:"domain"`
	}
	err := r.ExtractInto(&s)
	return s.Domain, err
}
