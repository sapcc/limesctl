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

package core

import (
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/sapcc/limes"
	"github.com/sapcc/limesctl/internal/errors"
)

type csvData [][]string

// renderCSV renders the result of a get/list/set operation in the CSV format.
func (c *Cluster) renderCSV() csvData {
	var data csvData
	var labels []string

	switch {
	case c.Output.Long:
		labels = []string{"cluster id", "area", "service", "category", "resource", "capacity",
			"domains quota", "usage", "physical usage", "burst usage", "unit", "comment", "scraped at (UTC)"}
	default:
		labels = []string{"cluster id", "service", "resource", "capacity", "domains quota", "usage", "unit"}
	}
	data = append(data, labels)

	if c.IsList {
		clusterList, err := c.Result.ExtractClusters()
		errors.Handle(err, "could not render the CSV data for clusters")

		for _, cluster := range clusterList {
			c.parseToCSV(&cluster, &data)
		}
	} else {
		cluster, err := c.Result.Extract()
		errors.Handle(err, "could not render the CSV data for cluster")

		c.parseToCSV(cluster, &data)
	}

	return data
}

// renderCSV renders the result of a get/list/set operation in the CSV format.
func (d *Domain) renderCSV() csvData {
	var data csvData
	var labels []string

	switch {
	case d.Output.Names:
		labels = []string{"domain name", "service", "resource", "quota", "projects quota", "usage", "unit"}
	case d.Output.Long:
		labels = []string{"domain id", "domain name", "area", "service", "category", "resource",
			"quota", "projects quota", "usage", "physical usage", "burst usage", "unit", "scraped at (UTC)"}
	default:
		labels = []string{"domain id", "service", "resource", "quota", "projects quota", "usage", "unit"}
	}
	data = append(data, labels)

	if d.IsList {
		domainList, err := d.Result.ExtractDomains()
		errors.Handle(err, "could not render the CSV data for domains")

		for _, domain := range domainList {
			d.parseToCSV(&domain, &data)
		}
	} else {
		domain, err := d.Result.Extract()
		errors.Handle(err, "could not render the CSV data for domain")

		d.parseToCSV(domain, &data)
	}

	return data
}

// renderCSV renders the result of a get/list/set operation in the CSV format.
func (p *Project) renderCSV() csvData {
	var data csvData
	var labels []string

	switch {
	case p.Output.Names:
		labels = []string{"domain name", "project name", "service", "resource", "quota", "usage", "unit"}
	case p.Output.Long:
		labels = []string{"domain id", "domain name", "project id", "project name",
			"area", "service", "category", "resource", "quota", "burst quota", "usage",
			"physical usage", "burst usage", "unit", "scraped at (UTC)"}
	default:
		labels = []string{"domain id", "project id", "service", "resource", "quota", "usage", "unit"}
	}
	data = append(data, labels)

	if p.IsList {
		projectList, err := p.Result.ExtractProjects()
		errors.Handle(err, "could not render the CSV data for projects")

		for _, project := range projectList {
			p.parseToCSV(&project, &data)
		}
	} else {
		project, err := p.Result.Extract()
		errors.Handle(err, "could not render the CSV data for project")

		p.parseToCSV(project, &data)
	}

	return data
}

// parseToCSV parses a limes.ClusterReport to CSV depending on the output format and assigns it
// to the aggregate csvData.
func (c *Cluster) parseToCSV(cluster *limes.ClusterReport, data *csvData) {
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
			if cSrvRes.Capacity != nil {
				cap = *cSrvRes.Capacity
			}
			physicalUsage := cSrvRes.Usage
			if cSrvRes.PhysicalUsage != nil {
				physicalUsage = *cSrvRes.PhysicalUsage
			}

			unit, val := humanReadable(c.Output.HumanReadable, cSrvRes.ResourceInfo.Unit, rawValues{
				"capacity":      cap,
				"domainsQuota":  cSrvRes.DomainsQuota,
				"usage":         cSrvRes.Usage,
				"burstUsage":    cSrvRes.BurstUsage,
				"physicalUsage": physicalUsage,
			})

			switch {
			case c.Output.Long:
				csvRecord = append(csvRecord, cluster.ID, cSrv.ServiceInfo.Area, cSrv.ServiceInfo.Type,
					cSrvRes.ResourceInfo.Category, cSrvRes.ResourceInfo.Name, val["capacity"], val["domainsQuota"],
					val["usage"], val["physicalUsage"], val["burstUsage"], unit, cSrvRes.Comment, timestampToString(cSrv.MinScrapedAt),
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

// parseToCSV parses a limes.DomainReport to CSV depending on the output format and assigns it
// to the aggregate csvData.
func (d *Domain) parseToCSV(domain *limes.DomainReport, data *csvData) {
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

			physicalUsage := dSrvRes.Usage
			if dSrvRes.PhysicalUsage != nil {
				physicalUsage = *dSrvRes.PhysicalUsage
			}

			unit, val := humanReadable(d.Output.HumanReadable, dSrvRes.ResourceInfo.Unit, rawValues{
				"domainQuota":   dSrvRes.DomainQuota,
				"projectsQuota": dSrvRes.ProjectsQuota,
				"usage":         dSrvRes.Usage,
				"burstUsage":    dSrvRes.BurstUsage,
				"physicalUsage": physicalUsage,
			})

			switch {
			case d.Output.Names:
				csvRecord = append(csvRecord, domain.Name, dSrv.ServiceInfo.Type, dSrvRes.ResourceInfo.Name,
					val["domainQuota"], val["projectsQuota"], val["usage"], unit,
				)
			case d.Output.Long:
				csvRecord = append(csvRecord, domain.UUID, domain.Name, dSrv.ServiceInfo.Area, dSrv.ServiceInfo.Type,
					dSrvRes.ResourceInfo.Category, dSrvRes.ResourceInfo.Name, val["domainQuota"], val["projectsQuota"],
					val["usage"], val["physicalUsage"], val["burstUsage"], unit, timestampToString(dSrv.MinScrapedAt),
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

// parseToCSV parses a limes.ProjectReport to CSV depending on the output format and assigns it
// to the aggregate csvData.
func (p *Project) parseToCSV(project *limes.ProjectReport, data *csvData) {
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

			var burstQuota, burstUsage uint64
			if project.Bursting != nil {
				if project.Bursting.Enabled {
					burstQuota = project.Bursting.Multiplier.ApplyTo(pSrvRes.Quota)

					if pSrvRes.Usage > pSrvRes.Quota {
						burstUsage = pSrvRes.Usage - pSrvRes.Quota
					}
				}
			}

			physicalUsage := pSrvRes.Usage
			if pSrvRes.PhysicalUsage != nil {
				physicalUsage = *pSrvRes.PhysicalUsage
			}

			unit, val := humanReadable(p.Output.HumanReadable, pSrvRes.ResourceInfo.Unit, rawValues{
				"quota":         pSrvRes.Quota,
				"burstQuota":    burstQuota,
				"usage":         pSrvRes.Usage,
				"burstUsage":    burstUsage,
				"physicalUsage": physicalUsage,
			})

			switch {
			case p.Output.Names:
				csvRecord = append(csvRecord, p.DomainName, project.Name, pSrv.ServiceInfo.Type,
					pSrvRes.ResourceInfo.Name, val["quota"], val["usage"], unit,
				)
			case p.Output.Long:
				csvRecord = append(csvRecord, p.DomainID, p.DomainName, project.UUID, project.Name, pSrv.ServiceInfo.Area,
					pSrv.ServiceInfo.Type, pSrvRes.ResourceInfo.Category, pSrvRes.ResourceInfo.Name, val["quota"],
					val["burstQuota"], val["usage"], val["physicalUsage"], val["burstUsage"], unit, timestampToString(pSrv.ScrapedAt),
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

func timestampToString(timestamp *int64) string {
	if timestamp == nil {
		return ""
	}
	return time.Unix(*timestamp, 0).UTC().Format(time.RFC3339)
}

type rawValues map[string]uint64
type convertedValues map[string]string

func humanReadable(convert bool, unit limes.Unit, rv rawValues) (string, convertedValues) {
	cv := make(convertedValues, len(rv))

	computeAgainst := smallestNonzeroValue(rv)

	if unit == limes.UnitNone || computeAgainst < 1024 {
		convert = false
	}
	if !convert {
		for k, v := range rv {
			cv[k] = strconv.FormatUint(v, 10)
		}
		return string(unit), cv
	}

	oldExp := quotaUnits[unit]
	usage := computeAgainst

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
	for k, v := range quotaUnits {
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

func smallestNonzeroValue(rv rawValues) uint64 {
	var vals []uint64
	for _, v := range rv {
		if v != 0 {
			vals = append(vals, v)
		}
	}
	if len(vals) == 0 {
		return 0
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	return vals[0]
}
