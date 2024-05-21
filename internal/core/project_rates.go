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

package core

import (
	"slices"

	"github.com/sapcc/go-api-declarations/limes"
	limesrates "github.com/sapcc/go-api-declarations/limes/rates"
)

// ProjectRatesReport is a wrapper for limesrates.ProjectReport.
type ProjectRatesReport struct {
	*limesrates.ProjectReport

	DomainID   string
	DomainName string
}

// LimesProjectRatesToReportRenderer wraps the given limesrates.ProjectReport in a
// ProjectRatesReport and returns a []LimesReportRenderer.
func LimesProjectRatesToReportRenderer(
	in []limesrates.ProjectReport,
	domainID, domainName string,
	hasRatesOnly bool,
) []LimesReportRenderer {

	out := make([]LimesReportRenderer, 0, len(in))
	for _, rep := range in {
		out = append(out, ProjectRatesReport{
			ProjectReport: &rep,
			DomainID:      domainID,
			DomainName:    domainName,
		})
	}
	return out
}

var csvHeaderProjectRatesDefault = []string{"domain id", "project id", "service", "rate", "limit", "window", "usage", "unit"}

var csvHeaderProjectRatesLong = []string{
	"domain id", "domain name", "project id", "project name", "area", "service",
	"rate", "limit", "default limit", "window", "default window", "usage", "unit", "scraped at (UTC)",
}

func (p ProjectRatesReport) getHeaderRow(opts *OutputOpts) []string {
	switch opts.CSVRecFmt {
	case CSVRecordFormatLong:
		return csvHeaderProjectRatesLong
	case CSVRecordFormatNames:
		h := csvHeaderProjectRatesDefault
		h[0] = domainName
		h[1] = "project name"
		return h
	default:
		return csvHeaderProjectRatesDefault
	}
}

func (p ProjectRatesReport) render(opts *OutputOpts) CSVRecords {
	var records CSVRecords

	// Serialize service types with ordered keys
	types := make([]limes.ServiceType, 0, len(p.Services))
	for typeStr := range p.Services {
		types = append(types, typeStr)
	}
	slices.Sort(types)

	for _, srv := range types {
		// Serialize rate names with ordered keys
		names := make([]limesrates.RateName, 0, len(p.Services[srv].Rates))
		for nameStr := range p.Services[srv].Rates {
			names = append(names, nameStr)
		}
		slices.Sort(names)

		for _, rate := range names {
			var r []string
			// Initialize temporary variables to make map lookups easier.
			pSrv := p.Services[srv]
			pSrvRate := p.Services[srv].Rates[rate]

			var window, defaultWindow string
			if pSrvRate.Window != nil {
				window = pSrvRate.Window.String()
			}
			if pSrvRate.DefaultWindow != nil {
				defaultWindow = pSrvRate.DefaultWindow.String()
			}

			valToStr := defaultValToStrFunc
			if opts.CSVRecFmt == CSVRecordFormatLong {
				r = append(r, p.DomainID, p.DomainName, p.UUID, p.Name, pSrv.Area, string(pSrv.Type), string(pSrvRate.Name),
					valToStr(pSrvRate.Limit), valToStr(pSrvRate.DefaultLimit), window, defaultWindow,
					pSrvRate.UsageAsBigint, string(pSrvRate.Unit), timestampToString(pSrv.ScrapedAt),
				)
			} else {
				projectNameOrID := p.UUID
				domainNameOrID := p.DomainID
				if opts.CSVRecFmt == CSVRecordFormatNames {
					projectNameOrID = p.Name
					domainNameOrID = p.DomainName
				}
				r = append(r, domainNameOrID, projectNameOrID, string(pSrv.Type), string(pSrvRate.Name),
					valToStr(pSrvRate.Limit), window, pSrvRate.UsageAsBigint, string(pSrvRate.Unit),
				)
			}

			records = append(records, r)
		}
	}

	return records
}
