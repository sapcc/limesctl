// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"slices"

	"github.com/sapcc/go-api-declarations/limes"
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
	types := make([]limes.ServiceType, 0, len(c.Services))
	for typeStr := range c.Services {
		types = append(types, typeStr)
	}
	slices.Sort(types)

	for _, srv := range types {
		// Serialize rate names with ordered keys
		names := make([]limesrates.RateName, 0, len(c.Services[srv].Rates))
		for nameStr := range c.Services[srv].Rates {
			names = append(names, nameStr)
		}
		slices.Sort(names)

		for _, rate := range names {
			var r []string
			// Initialize temporary variables to make map lookups easier.
			cSrv := c.Services[srv]
			cSrvRate := c.Services[srv].Rates[rate]

			valToStr := defaultValToStrFunc
			if opts.CSVRecFmt == CSVRecordFormatLong {
				r = append(r, c.ID, cSrv.Area, string(cSrv.Type), string(cSrvRate.Name), valToStr(cSrvRate.Limit),
					cSrvRate.Window.String(), string(cSrvRate.Unit), timestampToString(cSrv.MinScrapedAt))
			} else {
				r = append(r, c.ID, string(cSrv.Type), string(cSrvRate.Name), valToStr(cSrvRate.Limit), cSrvRate.Window.String(), string(cSrvRate.Unit))
			}

			records = append(records, r)
		}
	}

	return records
}
