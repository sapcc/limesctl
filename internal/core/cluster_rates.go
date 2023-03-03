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

	limesrates "github.com/sapcc/go-api-declarations/limes/rates"
)

// ClusterRatesReport is a wrapper for limesrates.ClusterReport.
type ClusterRatesReport struct {
	*limesrates.ClusterReport
}

var csvHeaderClusterRatesDefault = []string{"cluster id", "service", "rate", "limit", "window", "unit"}

var csvHeaderClusterRatesLong = []string{"cluster id", "area", "service", "rate", "limit",
	"window", "unit", "scraped at (UTC)"}

func (c ClusterRatesReport) getHeaderRow(opts *OutputOpts) []string {
	if opts.CSVRecFmt == CSVRecordFormatLong {
		return csvHeaderClusterRatesLong
	}
	return csvHeaderClusterRatesDefault
}

func (c ClusterRatesReport) render(opts *OutputOpts) CSVRecords {
	var records CSVRecords

	// Serialize service types with ordered keys
	types := make([]string, 0, len(c.Services))
	for typeStr := range c.Services {
		types = append(types, typeStr)
	}
	sort.Strings(types)

	for _, srv := range types {
		// Serialize rate names with ordered keys
		names := make([]string, 0, len(c.Services[srv].Rates))
		for nameStr := range c.Services[srv].Rates {
			names = append(names, nameStr)
		}
		sort.Strings(names)

		for _, rate := range names {
			var r []string
			// Initialize temporary variables to make map lookups easier.
			cSrv := c.Services[srv]
			cSrvRate := c.Services[srv].Rates[rate]

			valToStr := defaultValToStrFunc
			if opts.CSVRecFmt == CSVRecordFormatLong {
				r = append(r, c.ID, cSrv.Area, cSrv.Type, cSrvRate.Name, valToStr(cSrvRate.Limit),
					cSrvRate.Window.String(), string(cSrvRate.Unit), timestampToString(cSrv.MinScrapedAt))
			} else {
				r = append(r, c.ID, cSrv.Type, cSrvRate.Name, valToStr(cSrvRate.Limit), cSrvRate.Window.String(), string(cSrvRate.Unit))
			}

			records = append(records, r)
		}
	}

	return records
}
