// Copyright 2020 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
