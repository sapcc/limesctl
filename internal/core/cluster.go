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

// ClusterReport is a wrapper for limes.ClusterReport.
type ClusterReport struct {
	*limes.ClusterReport
}

var csvHeaderClusterLong = []string{
	"cluster id", "area", "service", "category", "resource", "capacity",
	"domains quota", "usage", "physical usage", "burst usage", "unit", "scraped at (UTC)",
}

var csvHeaderClusterDefault = []string{
	"cluster id", "service", "resource", "capacity", "domains quota",
	"usage", "unit",
}

// GetHeaderRow implements the LimesReportRenderer interface.
func (c ClusterReport) getHeaderRow(opts *OutputOpts) []string {
	if opts.CSVRecFmt == CSVRecordFormatLong {
		return csvHeaderClusterLong
	}
	return csvHeaderClusterDefault
}

// Render implements the LimesReportRenderer interface.
func (c ClusterReport) render(opts *OutputOpts) CSVRecords {
	var records CSVRecords

	// Serialize service types with ordered keys
	types := make([]string, 0, len(c.Services))
	for typeStr := range c.Services {
		types = append(types, typeStr)
	}
	sort.Strings(types)

	for _, srv := range types {
		// Serialize resource names with ordered keys
		names := make([]string, 0, len(c.Services[srv].Resources))
		for nameStr := range c.Services[srv].Resources {
			names = append(names, nameStr)
		}
		sort.Strings(names)

		for _, res := range names {
			var r []string
			// Initialize temporary variables to make map lookups easier.
			cSrv := c.Services[srv]
			cSrvRes := c.Services[srv].Resources[res]

			capacity := cSrvRes.Capacity
			physU := cSrvRes.PhysicalUsage
			domsQ := cSrvRes.DomainsQuota

			valToStr, unit := getValToStrFunc(opts.Humanize, cSrvRes.Unit, []uint64{
				zeroIfNil(capacity), zeroIfNil(physU), zeroIfNil(domsQ),
				cSrvRes.Usage, cSrvRes.BurstUsage,
			})

			if opts.CSVRecFmt == CSVRecordFormatLong {
				r = append(r, c.ID, cSrv.Area, cSrv.Type, cSrvRes.Category, cSrvRes.Name, emptyStrIfNil(capacity, valToStr),
					emptyStrIfNil(domsQ, valToStr), valToStr(cSrvRes.Usage), emptyStrIfNil(physU, valToStr),
					valToStr(cSrvRes.BurstUsage), string(unit), timestampToString(cSrv.MinScrapedAt),
				)
			} else {
				r = append(r, c.ID, cSrv.Type, cSrvRes.Name, emptyStrIfNil(capacity, valToStr),
					emptyStrIfNil(domsQ, valToStr), valToStr(cSrvRes.Usage), string(unit),
				)
			}

			records = append(records, r)
		}
	}

	return records
}
