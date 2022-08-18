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
	"sort"

	"github.com/sapcc/go-api-declarations/limes"
)

// ProjectReport is a wrapper for limes.ProjectReport.
type ProjectReport struct {
	*limes.ProjectReport

	HasRatesOnly bool

	DomainID   string
	DomainName string
}

// LimesProjectsToReportRenderer wraps the given limes.ProjectReport in a
// ProjectReport and returns a []LimesReportRenderer.
func LimesProjectsToReportRenderer(
	in []limes.ProjectReport,
	domainID, domainName string,
	hasRatesOnly bool,
) []LimesReportRenderer {

	out := make([]LimesReportRenderer, 0, len(in))
	for _, rep := range in {
		rep := rep
		out = append(out, ProjectReport{
			ProjectReport: &rep,
			DomainID:      domainID,
			DomainName:    domainName,
			HasRatesOnly:  hasRatesOnly,
		})
	}
	return out
}

var csvHeaderProjectDefault = []string{"domain id", "project id", "service", "resource", "quota", "usage", "unit"}

var csvHeaderProjectLong = []string{
	"domain id", "domain name", "project id", "project name", "area", "service",
	"category", "resource", "quota", "burst quota", "usage", "physical usage", "burst usage", "unit", "scraped at (UTC)",
}

// GetHeaderRow implements the LimesReportRenderer interface.
func (p ProjectReport) getHeaderRow(opts *OutputOpts) []string {
	if p.HasRatesOnly {
		return p.getRatesHeaderRow(opts)
	}

	switch opts.CSVRecFmt {
	case CSVRecordFormatLong:
		return csvHeaderProjectLong
	case CSVRecordFormatNames:
		h := csvHeaderProjectDefault
		h[0] = domainName
		h[1] = "project name"
		return h
	default:

		return csvHeaderProjectDefault
	}
}

var csvHeaderProjectRatesDefault = []string{"domain id", "project id", "service", "rate", "limit", "window", "usage", "unit"}

var csvHeaderProjectRatesLong = []string{
	"domain id", "domain name", "project id", "project name", "area", "service",
	"rate", "limit", "default limit", "window", "default window", "usage", "unit", "scraped at (UTC)",
}

func (p ProjectReport) getRatesHeaderRow(opts *OutputOpts) []string {
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

// Render implements the LimesReportRenderer interface.
func (p ProjectReport) render(opts *OutputOpts) CSVRecords {
	if p.HasRatesOnly {
		return p.renderRates(opts)
	}

	var records CSVRecords

	// Serialize service types with ordered keys
	types := make([]string, 0, len(p.Services))
	for typeStr := range p.Services {
		types = append(types, typeStr)
	}
	sort.Strings(types)

	for _, srv := range types {
		// Serialize resource names with ordered keys
		names := make([]string, 0, len(p.Services[srv].Resources))
		for nameStr := range p.Services[srv].Resources {
			names = append(names, nameStr)
		}
		sort.Strings(names)

		for _, res := range names {
			var r []string
			// Initialize temporary variables to make map lookups easier.
			pSrv := p.Services[srv]
			pSrvRes := p.Services[srv].Resources[res]

			physU := pSrvRes.PhysicalUsage
			quota := pSrvRes.Quota
			usage := pSrvRes.Usage

			// We use a *uint64 for burstQuota instead of an uint64 for
			// consistency with Limes' API, i.e. if quota has a null value
			// then burstQuota should also be null instead of zero.
			var burstQuota *uint64
			var burstUsage uint64
			if quota != nil && p.Bursting != nil && p.Bursting.Enabled {
				q := *quota
				bq := p.Bursting.Multiplier.ApplyTo(q)
				burstQuota = &bq
				if usage > q {
					burstUsage = usage - q
				}
			}

			valToStr, unit := getValToStrFunc(opts.Humanize, pSrvRes.Unit, []uint64{
				zeroIfNil(burstQuota), burstUsage, zeroIfNil(physU), zeroIfNil(quota), usage,
			})

			if opts.CSVRecFmt == CSVRecordFormatLong {
				r = append(r, p.DomainID, p.DomainName, p.UUID, p.Name, pSrv.Area, pSrv.Type, pSrvRes.Category,
					pSrvRes.Name, emptyStrIfNil(quota, valToStr), emptyStrIfNil(burstQuota, valToStr), valToStr(usage),
					emptyStrIfNil(physU, valToStr), valToStr(burstUsage), string(unit), timestampToString(pSrv.ScrapedAt),
				)
			} else {
				projectNameOrID := p.UUID
				domainNameOrID := p.DomainID
				if opts.CSVRecFmt == CSVRecordFormatNames {
					projectNameOrID = p.Name
					domainNameOrID = p.DomainName
				}
				r = append(r, domainNameOrID, projectNameOrID, pSrv.Type, pSrvRes.Name,
					emptyStrIfNil(quota, valToStr), valToStr(usage), string(unit),
				)
			}

			records = append(records, r)
		}
	}

	return records
}

func (p ProjectReport) renderRates(opts *OutputOpts) CSVRecords {
	var records CSVRecords

	// Serialize service types with ordered keys
	types := make([]string, 0, len(p.Services))
	for typeStr := range p.Services {
		types = append(types, typeStr)
	}
	sort.Strings(types)

	for _, srv := range types {
		// Serialize rate names with ordered keys
		names := make([]string, 0, len(p.Services[srv].Rates))
		for nameStr := range p.Services[srv].Rates {
			names = append(names, nameStr)
		}
		sort.Strings(names)

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
				r = append(r, p.DomainID, p.DomainName, p.UUID, p.Name, pSrv.Area, pSrv.Type, pSrvRate.Name,
					valToStr(pSrvRate.Limit), valToStr(pSrvRate.DefaultLimit), window, defaultWindow,
					pSrvRate.UsageAsBigint, string(pSrvRate.Unit), timestampToString(pSrv.RatesScrapedAt),
				)
			} else {
				projectNameOrID := p.UUID
				domainNameOrID := p.DomainID
				if opts.CSVRecFmt == CSVRecordFormatNames {
					projectNameOrID = p.Name
					domainNameOrID = p.DomainName
				}
				r = append(r, domainNameOrID, projectNameOrID, pSrv.Type, pSrvRate.Name,
					valToStr(pSrvRate.Limit), window, pSrvRate.UsageAsBigint, string(pSrvRate.Unit),
				)
			}

			records = append(records, r)
		}
	}

	return records
}
