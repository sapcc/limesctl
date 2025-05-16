// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"slices"

	"github.com/sapcc/go-api-declarations/limes"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
)

// DomainReport is a wrapper for limesresources.DomainReport.
type DomainReport struct {
	*limesresources.DomainReport
}

// LimesDomainsToReportRenderer wraps the given limesresources.DomainReport in a
// DomainReport and returns a []LimesReportRenderer.
func LimesDomainsToReportRenderer(in []limesresources.DomainReport) []LimesReportRenderer {
	out := make([]LimesReportRenderer, 0, len(in))
	for _, rep := range in {
		out = append(out, DomainReport{DomainReport: &rep})
	}
	return out
}

var csvHeaderDomainDefault = []string{"domain id", "service", "resource", "quota", "projects quota", "usage", "unit"}

var csvHeaderDomainLong = []string{
	"domain id", "domain name", "area", "service", "category", "resource",
	"quota", "projects quota", "usage", "physical usage", "unit", "scraped at (UTC)",
}

const domainName = "domain name"

// GetHeaderRow implements the LimesReportRenderer interface.
func (d DomainReport) getHeaderRow(opts *OutputOpts) []string {
	switch opts.CSVRecFmt {
	case CSVRecordFormatLong:
		return csvHeaderDomainLong
	case CSVRecordFormatNames:
		h := csvHeaderDomainDefault
		h[0] = domainName
		return h
	default:
		return csvHeaderDomainDefault
	}
}

// Render implements the LimesReportRenderer interface.
func (d DomainReport) render(opts *OutputOpts) CSVRecords {
	var records CSVRecords

	// Serialize service types with ordered keys
	types := make([]limes.ServiceType, 0, len(d.Services))
	for typeStr := range d.Services {
		types = append(types, typeStr)
	}
	slices.Sort(types)

	for _, srv := range types {
		// Serialize resource names with ordered keys
		names := make([]limesresources.ResourceName, 0, len(d.Services[srv].Resources))
		for nameStr := range d.Services[srv].Resources {
			names = append(names, nameStr)
		}
		slices.Sort(names)

		for _, res := range names {
			var r []string
			// Initialize temporary variables to make map lookups easier.
			dSrv := d.Services[srv]
			dSrvRes := d.Services[srv].Resources[res]

			physU := dSrvRes.PhysicalUsage
			domQ := dSrvRes.DomainQuota
			projectsQ := dSrvRes.ProjectsQuota

			valToStr, unit := getValToStrFunc(opts.Humanize, dSrvRes.Unit, []uint64{
				zeroIfNil(physU), zeroIfNil(domQ), zeroIfNil(projectsQ),
				dSrvRes.Usage,
			})

			if opts.CSVRecFmt == CSVRecordFormatLong {
				r = append(r, d.UUID, d.Name, dSrv.Area, string(dSrv.Type), dSrvRes.Category, string(dSrvRes.Name),
					emptyStrIfNil(domQ, valToStr), emptyStrIfNil(projectsQ, valToStr), valToStr(dSrvRes.Usage),
					emptyStrIfNil(physU, valToStr), string(unit), timestampToString(dSrv.MinScrapedAt),
				)
			} else {
				nameOrID := d.UUID
				if opts.CSVRecFmt == CSVRecordFormatNames {
					nameOrID = d.Name
				}
				r = append(r, nameOrID, string(dSrv.Type), string(dSrvRes.Name), emptyStrIfNil(domQ, valToStr),
					emptyStrIfNil(projectsQ, valToStr), valToStr(dSrvRes.Usage), string(unit),
				)
			}

			records = append(records, r)
		}
	}

	return records
}
