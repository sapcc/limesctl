/*******************************************************************************
*
* Copyright 2018 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package cli

import (
	"errors"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/sapcc/limes/pkg/limes"
	"github.com/sapcc/limes/pkg/reports"
)

type csvData [][]string

// renderCSV renders the result of a get/list/set operation in the CSV format.
func (c *Cluster) renderCSV() *csvData {
	var data csvData
	var labels []string

	switch {
	case c.Opts.Long:
		labels = []string{"cluster id", "area", "service", "category", "resource", "capacity",
			"domains quota", "usage", "unit", "comment", "scraped at (UTC)"}
	default:
		labels = []string{"cluster id", "service", "resource", "capacity", "domains quota", "usage", "unit"}
	}

	if c.IsList {
		clusterList, err := c.Result.ExtractClusters()
		handleError("could not render the CSV data for clusters", err)

		data = append(data, labels)
		for _, cluster := range clusterList {
			c.parseToCSV(&cluster, &data)
		}
	} else {
		cluster, err := c.Result.Extract()
		handleError("could not render the CSV data for cluster", err)

		data = append(data, labels)
		c.parseToCSV(cluster, &data)
	}

	return &data
}

// renderCSV renders the result of a get/list/set operation in the CSV format.
func (d *Domain) renderCSV() *csvData {
	var data csvData
	var labels []string

	if d.Opts.Names && d.Opts.Long {
		handleError("", errors.New("'--names' and '--long' can not be used together"))
	}

	switch {
	case d.Opts.Names:
		labels = []string{"domain name", "service", "resource", "quota", "projects quota", "usage", "unit"}
	case d.Opts.Long:
		labels = []string{"domain id", "domain name", "area", "service", "category", "resource",
			"quota", "projects quota", "usage", "unit", "scraped at (UTC)"}
	default:
		labels = []string{"domain id", "service", "resource", "quota", "projects quota", "usage", "unit"}
	}

	if d.IsList {
		domainList, err := d.Result.ExtractDomains()
		handleError("could not render the CSV data for domains", err)

		data = append(data, labels)
		for _, domain := range domainList {
			d.parseToCSV(&domain, &data)
		}
	} else {
		domain, err := d.Result.Extract()
		handleError("could not render the CSV data for domain", err)

		data = append(data, labels)
		d.parseToCSV(domain, &data)
	}

	return &data
}

// renderCSV renders the result of a get/list/set operation in the CSV format.
func (p *Project) renderCSV() *csvData {
	var data csvData
	var labels []string

	if p.Opts.Names && p.Opts.Long {
		handleError("", errors.New("'--names' and '--long' can not be used together"))
	}

	switch {
	case p.Opts.Names:
		labels = []string{"domain name", "project name", "service", "resource", "quota", "usage", "unit"}
	case p.Opts.Long:
		labels = []string{"domain id", "domain name", "project id", "project name", "area",
			"service", "category", "resource", "quota", "usage", "unit", "scraped at (UTC)"}
	default:
		labels = []string{"domain id", "project id", "service", "resource", "quota", "usage", "unit"}
	}

	if p.IsList {
		projectList, err := p.Result.ExtractProjects()
		handleError("could not render the CSV data for projects", err)

		data = append(data, labels)
		for _, project := range projectList {
			p.parseToCSV(&project, &data)
		}
	} else {
		project, err := p.Result.Extract()
		handleError("could not render the CSV data for project", err)

		data = append(data, labels)
		p.parseToCSV(project, &data)
	}

	return &data
}

// parseToCSV parses a reports.Cluster to CSV depending on the output format and assigns it
// to the aggregate csvData.
func (c *Cluster) parseToCSV(cluster *reports.Cluster, data *csvData) {
	//serialize service types with ordered keys
	types := make([]string, 0, len(cluster.Services))
	for typeStr := range cluster.Services {
		types = append(types, typeStr)
	}
	sort.Strings(types)

	for _, srv := range types {
		//serialize resource names with ordered keys
		names := make([]string, 0, len(cluster.Services[srv].Resources))
		for nameStr := range cluster.Services[srv].Resources {
			names = append(names, nameStr)
		}
		sort.Strings(names)

		for _, res := range names {
			var csvRecord []string
			// temporary variables to make map lookups easier
			cSrv := cluster.Services[srv]
			cSrvRes := cluster.Services[srv].Resources[res]

			// need to do this check to avoid nil pointers
			var cap uint64
			if tmp := cSrvRes.Capacity; tmp != nil {
				cap = *tmp
			}

			unit, val := humanReadable(c.Opts.HumanReadable, cSrvRes.ResourceInfo.Unit, rawValues{
				"capacity":     cap,
				"domainsQuota": cSrvRes.DomainsQuota,
				"usage":        cSrvRes.Usage,
			})

			switch {
			case c.Opts.Long:
				csvRecord = append(csvRecord, cluster.ID, cSrv.ServiceInfo.Area, cSrv.ServiceInfo.Type,
					cSrvRes.ResourceInfo.Category, cSrvRes.ResourceInfo.Name, val["capacity"], val["domainsQuota"],
					val["usage"], unit, cSrvRes.Comment, time.Unix(cSrv.MinScrapedAt, 0).UTC().Format(time.RFC3339),
				)
			default:
				csvRecord = append(csvRecord, cluster.ID, cSrv.ServiceInfo.Type, cSrvRes.ResourceInfo.Name,
					val["capacity"], val["domainsQuota"], val["usage"], unit,
				)
			}

			*data = append(*data, csvRecord)
		}
	}
}

// parseToCSV parses a reports.Domain to CSV depending on the output format and assigns it
// to the aggregate csvData.
func (d *Domain) parseToCSV(domain *reports.Domain, data *csvData) {
	//serialize service types with ordered keys
	types := make([]string, 0, len(domain.Services))
	for typeStr := range domain.Services {
		types = append(types, typeStr)
	}
	sort.Strings(types)

	for _, srv := range types {
		//serialize resource names with ordered keys
		names := make([]string, 0, len(domain.Services[srv].Resources))
		for nameStr := range domain.Services[srv].Resources {
			names = append(names, nameStr)
		}
		sort.Strings(names)

		for _, res := range names {
			var csvRecord []string
			// temporary variables to make map lookups easier
			dSrv := domain.Services[srv]
			dSrvRes := domain.Services[srv].Resources[res]

			unit, val := humanReadable(d.Opts.HumanReadable, dSrvRes.ResourceInfo.Unit, rawValues{
				"domainQuota":   dSrvRes.DomainQuota,
				"projectsQuota": dSrvRes.ProjectsQuota,
				"usage":         dSrvRes.Usage,
			})

			switch {
			case d.Opts.Names:
				csvRecord = append(csvRecord, domain.Name, dSrv.ServiceInfo.Type, dSrvRes.ResourceInfo.Name,
					val["domainQuota"], val["projectsQuota"], val["usage"], unit,
				)
			case d.Opts.Long:
				csvRecord = append(csvRecord, domain.UUID, domain.Name, dSrv.ServiceInfo.Area, dSrv.ServiceInfo.Type,
					dSrvRes.ResourceInfo.Category, dSrvRes.ResourceInfo.Name, val["domainQuota"], val["projectsQuota"],
					val["usage"], unit, time.Unix(dSrv.MinScrapedAt, 0).UTC().Format(time.RFC3339),
				)
			default:
				csvRecord = append(csvRecord, domain.UUID, dSrv.ServiceInfo.Type, dSrvRes.ResourceInfo.Name,
					val["domainQuota"], val["projectsQuota"], val["usage"], unit,
				)
			}

			*data = append(*data, csvRecord)
		}
	}
}

// parseToCSV parses a reports.Project to CSV depending on the output format and assigns it
// to the aggregate csvData.
func (p *Project) parseToCSV(project *reports.Project, data *csvData) {
	//serialize service types with ordered keys
	types := make([]string, 0, len(project.Services))
	for typeStr := range project.Services {
		types = append(types, typeStr)
	}
	sort.Strings(types)

	for _, srv := range types {
		//serialize resource names with ordered keys
		names := make([]string, 0, len(project.Services[srv].Resources))
		for nameStr := range project.Services[srv].Resources {
			names = append(names, nameStr)
		}
		sort.Strings(names)

		for _, res := range names {
			var csvRecord []string
			// temporary variables to make map lookups easier
			pSrv := project.Services[srv]
			pSrvRes := project.Services[srv].Resources[res]

			unit, val := humanReadable(p.Opts.HumanReadable, pSrvRes.ResourceInfo.Unit, rawValues{
				"quota": pSrvRes.Quota,
				"usage": pSrvRes.Usage,
			})

			switch {
			case p.Opts.Names:
				csvRecord = append(csvRecord, p.DomainName, project.Name, pSrv.ServiceInfo.Type,
					pSrvRes.ResourceInfo.Name, val["quota"], val["usage"], unit,
				)
			case p.Opts.Long:
				csvRecord = append(csvRecord, p.DomainID, p.DomainName, project.UUID, project.Name, pSrv.ServiceInfo.Area,
					pSrv.ServiceInfo.Type, pSrvRes.ResourceInfo.Category, pSrvRes.ResourceInfo.Name,
					val["quota"], val["usage"], unit, time.Unix(pSrv.ScrapedAt, 0).UTC().Format(time.RFC3339),
				)
			default:
				csvRecord = append(csvRecord, p.DomainID, project.UUID, pSrv.ServiceInfo.Type,
					pSrvRes.ResourceInfo.Name, val["quota"], val["usage"], unit,
				)
			}

			*data = append(*data, csvRecord)
		}
	}
}

type rawValues map[string]uint64
type convertedValues map[string]string

// unitExp is map of limes.Unit to x, such that base-2 exponential of x
// is the number of bytes for some specific unit.
var unitExp = map[limes.Unit]float64{
	limes.UnitKibibytes: 10,
	limes.UnitMebibytes: 20,
	limes.UnitGibibytes: 30,
	limes.UnitTebibytes: 40,
	limes.UnitPebibytes: 50,
	limes.UnitExbibytes: 60,
}

func humanReadable(convert bool, unit limes.Unit, rv rawValues) (string, convertedValues) {
	cv := make(convertedValues, len(rv))

	// 'usage' is used to determine conversion suitability because it is common
	// in all three hierarchies and is the lowest value amongst its counterparts
	if unit == limes.UnitNone || rv["usage"] < 1024 {
		convert = false
	}
	if !convert {
		for k, v := range rv {
			cv[k] = strconv.FormatUint(v, 10)
		}
		return string(unit), cv
	}

	oldExp := unitExp[unit]
	usage := rv["usage"]

	var diffInExp float64
	// 2^60 bytes (exbibytes) is the maximum supported unit
	for diffInExp = 10; diffInExp <= (60 - oldExp); diffInExp += 10 {
		usageScaled := usage / uint64(math.Exp2(diffInExp))

		if usageScaled < 1024 {
			break
		}
	}

	// determine the new unit
	var newUnit limes.Unit
	for k, v := range unitExp {
		if v == (oldExp + diffInExp) {
			newUnit = k
		}
	}

	// convert values to the new unit
	for k, v := range rv {
		v := float64(v)
		v = v / math.Exp2(diffInExp)
		// round to second decimal place
		v = math.Round(v*100) / 100
		cv[k] = strconv.FormatFloat(v, 'f', -1, 64)
	}

	return string(newUnit), cv
}
