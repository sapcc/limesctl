package core

import (
	"sort"

	"github.com/sapcc/limes"
)

// DomainReport is a wrapper for limes.DomainReport.
type DomainReport struct {
	*limes.DomainReport
}

var csvHeaderDomainDefault = []string{"domain id", "service", "resource", "quota", "projects quota", "usage", "unit"}

var csvHeaderDomainLong = []string{"domain id", "domain name", "area", "service", "category", "resource",
	"quota", "projects quota", "usage", "physical usage", "burst usage", "unit", "scraped at (UTC)"}

// GetHeaderRow implements the LimesReportRenderer interface.
func (d DomainReport) getHeaderRow(csvFmt CSVRecordFormat) []string {
	switch csvFmt {
	case CSVRecordFormatLong:
		return csvHeaderDomainLong
	case CSVRecordFormatNames:
		h := csvHeaderDomainDefault
		h[0] = "domain name"
		return h
	default:
		return csvHeaderDomainDefault
	}
}

// Render implements the LimesReportRenderer interface.
func (d DomainReport) render(csvFmt CSVRecordFormat, humanize bool) CSVRecords {
	var records CSVRecords

	// Serialize service types with ordered keys
	types := make([]string, 0, len(d.Services))
	for typeStr := range d.Services {
		types = append(types, typeStr)
	}
	sort.Strings(types)

	for _, srv := range types {
		// Serialize resource names with ordered keys
		names := make([]string, 0, len(d.Services[srv].Resources))
		for nameStr := range d.Services[srv].Resources {
			names = append(names, nameStr)
		}
		sort.Strings(names)

		for _, res := range names {
			var r []string
			// Initialize temporary variables to make map lookups easier.
			dSrv := d.Services[srv]
			dSrvRes := d.Services[srv].Resources[res]

			// This check is necessary to avoid nil pointer dereference.
			physicalUsage := valFromPtr(dSrvRes.PhysicalUsage)
			domainQuota := valFromPtr(dSrvRes.DomainQuota)
			projectsQuota := valFromPtr(dSrvRes.ProjectsQuota)

			valToStr, unit := getValToStrFunc(humanize, dSrvRes.Unit, []uint64{
				physicalUsage, domainQuota, projectsQuota,
				dSrvRes.Usage, dSrvRes.BurstUsage,
			})

			physicalUsageStr := emptyStrIfZero(valToStr(physicalUsage))

			if csvFmt == CSVRecordFormatLong {
				r = append(r, d.UUID, d.Name, dSrv.Area, dSrv.Type, dSrvRes.Category, dSrvRes.Name,
					valToStr(domainQuota), valToStr(projectsQuota), valToStr(dSrvRes.Usage),
					physicalUsageStr, valToStr(dSrvRes.BurstUsage), string(unit), timestampToString(dSrv.MinScrapedAt),
				)
			} else {
				nameOrID := d.UUID
				if csvFmt == CSVRecordFormatNames {
					nameOrID = d.Name
				}
				r = append(r, nameOrID, dSrv.Type, dSrvRes.Name, valToStr(domainQuota),
					valToStr(projectsQuota), valToStr(dSrvRes.Usage), string(unit),
				)
			}

			records = append(records, r)
		}
	}

	return records
}
