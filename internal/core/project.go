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

	"github.com/sapcc/limes"
)

// ProjectReport is a wrapper for limes.ProjectReport.
type ProjectReport struct {
	*limes.ProjectReport

	DomainID   string
	DomainName string
}

// LimesProjectsToReportRenderer wraps the given limes.ProjectReport in a
// ProjectReport and returns a []LimesReportRenderer.
func LimesProjectsToReportRenderer(in []limes.ProjectReport, domainID, domainName string) []LimesReportRenderer {
	out := make([]LimesReportRenderer, 0, len(in))
	for _, rep := range in {
		rep := rep
		out = append(out, ProjectReport{
			ProjectReport: &rep,
			DomainID:      domainID,
			DomainName:    domainName,
		})
	}
	return out
}

var csvHeaderProjectDefault = []string{"domain id", "project id", "service", "resource", "quota", "usage", "unit"}

var csvHeaderProjectLong = []string{"domain id", "domain name", "project id", "project name", "area", "service",
	"category", "resource", "quota", "burst quota", "usage", "physical usage", "burst usage", "unit", "scraped at (UTC)"}

// GetHeaderRow implements the LimesReportRenderer interface.
func (p ProjectReport) getHeaderRow(csvFmt CSVRecordFormat) []string {
	switch csvFmt {
	case CSVRecordFormatLong:
		return csvHeaderProjectLong
	case CSVRecordFormatNames:
		h := csvHeaderProjectDefault
		h[0] = "domain name"
		h[1] = "project name"
		return h
	default:
		return csvHeaderProjectDefault
	}
}

// Render implements the LimesReportRenderer interface.
func (p ProjectReport) render(csvFmt CSVRecordFormat, humanize bool) CSVRecords {
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

			physicalUsage := valFromPtr(pSrvRes.PhysicalUsage)
			quota := valFromPtr(pSrvRes.Quota)
			usage := pSrvRes.Usage

			var burstQuota, burstUsage uint64
			if p.Bursting != nil && p.Bursting.Enabled {
				burstQuota = p.Bursting.Multiplier.ApplyTo(quota)
				if usage > quota {
					burstUsage = usage - quota
				}
			}

			valToStr, unit := getValToStrFunc(humanize, pSrvRes.Unit, []uint64{
				burstQuota, burstUsage, physicalUsage, quota, usage,
			})

			physicalUsageStr := emptyStrIfZero(valToStr(physicalUsage))
			quotaStr := emptyStrIfZero(valToStr(quota))
			burstQuotaStr := emptyStrIfZero(valToStr(burstQuota))

			if csvFmt == CSVRecordFormatLong {
				r = append(r, p.DomainID, p.DomainName, p.UUID, p.Name, pSrv.Area, pSrv.Type, pSrvRes.Category,
					pSrvRes.Name, quotaStr, burstQuotaStr, valToStr(usage),
					physicalUsageStr, valToStr(burstUsage), string(unit), timestampToString(pSrv.ScrapedAt),
				)
			} else {
				projectNameOrID := p.UUID
				domainNameOrID := p.DomainID
				if csvFmt == CSVRecordFormatNames {
					projectNameOrID = p.Name
					domainNameOrID = p.DomainName
				}
				r = append(r, domainNameOrID, projectNameOrID, pSrv.Type, pSrvRes.Name,
					quotaStr, valToStr(usage), string(unit),
				)
			}

			records = append(records, r)
		}
	}

	return records
}
