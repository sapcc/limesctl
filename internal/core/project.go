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

	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
)

// ProjectResourcesReport is a wrapper for limesresources.ProjectReport.
type ProjectResourcesReport struct {
	*limesresources.ProjectReport

	DomainID   string
	DomainName string
}

// LimesProjectResourcesToReportRenderer wraps the given limesresources.ProjectReport in a
// ProjectResourcesReport and returns a []LimesReportRenderer.
func LimesProjectResourcesToReportRenderer(
	in []limesresources.ProjectReport,
	domainID, domainName string,
	hasRatesOnly bool,
) []LimesReportRenderer {

	out := make([]LimesReportRenderer, 0, len(in))
	for _, rep := range in {
		rep := rep
		out = append(out, ProjectResourcesReport{
			ProjectReport: &rep,
			DomainID:      domainID,
			DomainName:    domainName,
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
func (p ProjectResourcesReport) getHeaderRow(opts *OutputOpts) []string {
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

// Render implements the LimesReportRenderer interface.
func (p ProjectResourcesReport) render(opts *OutputOpts) CSVRecords {
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
				bq := p.Bursting.Multiplier.ApplyTo(q, pSrvRes.QuotaDistributionModel)
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
