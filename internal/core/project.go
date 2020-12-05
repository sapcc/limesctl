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

			var burstQuota, burstUsage uint64
			if p.Bursting != nil && p.Bursting.Enabled {
				burstQuota = p.Bursting.Multiplier.ApplyTo(quota)
				if pSrvRes.Usage > quota {
					burstUsage = pSrvRes.Usage - quota
				}
			}

			valToStr, unit := getValToStrFunc(humanize, pSrvRes.Unit, []uint64{
				burstQuota, burstUsage, physicalUsage, quota, pSrvRes.Usage,
			})

			physicalUsageStr := emptyStrIfZero(valToStr(physicalUsage))

			if csvFmt == CSVRecordFormatLong {
				r = append(r, p.DomainID, p.DomainName, p.UUID, p.Name, pSrv.Area, pSrv.Type, pSrvRes.Category,
					pSrvRes.Name, valToStr(quota), valToStr(burstQuota), valToStr(pSrvRes.Usage),
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
					valToStr(quota), valToStr(pSrvRes.Usage), string(unit),
				)
			}

			records = append(records, r)
		}
	}

	return records
}
