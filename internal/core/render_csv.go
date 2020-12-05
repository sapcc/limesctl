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
	"sort"

	"github.com/sapcc/limes"
	"github.com/sapcc/limesctl/internal/errors"
)

type csvData [][]string

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

			var burstQuota, burstUsage, physicalUsage uint64
			if project.Bursting != nil {
				if project.Bursting.Enabled {
					burstQuota = project.Bursting.Multiplier.ApplyTo(pSrvRes.Quota)

					if pSrvRes.Usage > pSrvRes.Quota {
						burstUsage = pSrvRes.Usage - pSrvRes.Quota
					}
				}
			}
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
			if val["physicalUsage"] == "0" {
				val["physicalUsage"] = ""
			}

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
