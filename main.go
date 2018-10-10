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

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/sapcc/limes/pkg/api"
	"github.com/sapcc/limes/pkg/limes"
	"github.com/sapcc/limesctl/pkg/cli"
)

var (
	version = "1.0.0"
	app     = kingpin.New("limesctl", "CLI client for Limes.")

	// first-level commands and flags
	clusterCmd = app.Command("cluster", "Do some action on cluster(s).")
	domainCmd  = app.Command("domain", "Do some action on domain(s).")
	projectCmd = app.Command("project", "Do some action on project(s).")

	area        = app.Flag("area", "Resource area.").String()
	service     = app.Flag("service", "Service type").String()
	resource    = app.Flag("resource", "Resource name.").String()
	namesOutput = app.Flag("names", "Show output with names instead of UUIDs.").Bool()
	longOutput  = app.Flag("long", "Show detailed output.").Bool()
	outputFmt   = app.Flag("format", "Output format (table, json, csv).").PlaceHolder("table").Short('f').Enum("table", "json", "csv")

	// second-level subcommands and their flags/args
	clusterListCmd = clusterCmd.Command("list", "Query data for all the clusters. Requires a cloud-admin token.")

	clusterShowCmd = clusterCmd.Command("show", "Query data for a specific cluster, defaults to current cluster. Requires a cloud-admin token.")
	clusterShowID  = clusterShowCmd.Arg("cluster-id", "Cluster ID.").Required().String()

	clusterSetCmd  = clusterCmd.Command("set", "Change resource(s) quota for a cluster. Requires a cloud-admin token.")
	clusterSetID   = clusterSetCmd.Arg("cluster-id", "Cluster ID.").Required().String()
	clusterSetCaps = Capacity(clusterSetCmd.Arg("quota", "Quota value(s) to change. Format: service/resource=value(unit)").Required())

	domainListCmd = domainCmd.Command("list", "Query data for all the domains. Requires a cloud-admin token.")

	domainShowCmd = domainCmd.Command("show", "Query data for a specific domain. Requires a domain-admin token.")
	domainShowID  = domainShowCmd.Arg("domain-id", "Domain ID (name/UUID).").Required().String()

	domainSetCmd    = domainCmd.Command("set", "Change resource(s) quota for a domain. Requires a cloud-admin token.")
	domainSetID     = domainSetCmd.Arg("domain-id", "Domain ID (name/UUID).").Required().String()
	domainSetQuotas = Quota(domainSetCmd.Arg("quota", "Quota value(s) to change. Format: service/resource=value(unit)").Required())

	projectListCmd    = projectCmd.Command("list", "Query data for all the projects in a domain. Requires a domain-admin token.")
	projectListDomain = projectListCmd.Flag("domain", "Domain ID.").Short('d').Required().String()

	projectShowCmd    = projectCmd.Command("show", "Query data for a specific project in a domain. Requires project member permissions.")
	projectShowDomain = projectShowCmd.Flag("domain", "Domain ID.").Short('d').String()
	projectShowID     = projectShowCmd.Arg("project-id", "Project ID (name/UUID).").Required().String()

	projectSetCmd    = projectCmd.Command("set", "Change resource(s) quota for a project. Requires a domain-admin token.")
	projectSetDomain = projectSetCmd.Flag("domain", "Domain ID.").Short('d').String()
	projectSetID     = projectSetCmd.Arg("project-id", "Project ID (name/UUID).").Required().String()
	projectSetQuotas = Quota(projectSetCmd.Arg("quota", "Quota value(s) to change. Format: service/resource=value(unit)").Required())

	projectSyncCmd    = projectCmd.Command("sync", "Sync a project's quota and usage data from the backing services into Limes' local database. Requires a project-admin token.")
	projectSyncDomain = projectSyncCmd.Flag("domain", "Domain ID.").Short('d').String()
	projectSyncID     = projectSyncCmd.Arg("project-id", "Project ID (name/UUID).").Required().String()
)

func main() {
	app.Version("limesctl version " + version)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case clusterListCmd.FullCommand():
		c := cli.Cluster{
			Opts: cli.Options{
				Long:     *longOutput,
				Area:     *area,
				Service:  *service,
				Resource: *resource,
			},
		}
		cli.RunListTask(&c, *outputFmt)
	case clusterShowCmd.FullCommand():
		c := cli.Cluster{
			ID: *clusterShowID,
			Opts: cli.Options{
				Long:     *longOutput,
				Area:     *area,
				Service:  *service,
				Resource: *resource,
			},
		}
		cli.RunGetTask(&c, *outputFmt)
	case clusterSetCmd.FullCommand():
		c := cli.Cluster{
			ID: *clusterSetID,
			Opts: cli.Options{
				Names: *namesOutput,
				Long:  *longOutput,
			},
		}
		cli.RunSetTask(&c, *clusterSetCaps, *outputFmt)

	case domainListCmd.FullCommand():
		d := cli.Domain{
			Opts: cli.Options{
				Names:    *namesOutput,
				Long:     *longOutput,
				Area:     *area,
				Service:  *service,
				Resource: *resource,
			},
		}
		cli.RunListTask(&d, *outputFmt)
	case domainShowCmd.FullCommand():
		d := cli.Domain{
			ID: *domainShowID,
			Opts: cli.Options{
				Names:    *namesOutput,
				Long:     *longOutput,
				Area:     *area,
				Service:  *service,
				Resource: *resource,
			},
		}
		cli.RunGetTask(&d, *outputFmt)
	case domainSetCmd.FullCommand():
		d := cli.Domain{
			ID: *domainSetID,
			Opts: cli.Options{
				Names: *namesOutput,
				Long:  *longOutput,
			},
		}
		cli.RunSetTask(&d, *domainSetQuotas, *outputFmt)

	case projectListCmd.FullCommand():
		p := cli.Project{
			DomainID: *projectListDomain,
			Opts: cli.Options{
				Names:    *namesOutput,
				Long:     *longOutput,
				Area:     *area,
				Service:  *service,
				Resource: *resource,
			},
		}
		cli.RunListTask(&p, *outputFmt)
	case projectShowCmd.FullCommand():
		p := cli.Project{
			ID:       *projectShowID,
			DomainID: *projectShowDomain,
			Opts: cli.Options{
				Names:    *namesOutput,
				Long:     *longOutput,
				Area:     *area,
				Service:  *service,
				Resource: *resource,
			},
		}
		cli.RunGetTask(&p, *outputFmt)
	case projectSetCmd.FullCommand():
		p := cli.Project{
			ID:       *projectSetID,
			DomainID: *projectSetDomain,
			Opts: cli.Options{
				Names: *namesOutput,
				Long:  *longOutput,
			},
		}
		cli.RunSetTask(&p, *projectSetQuotas, *outputFmt)
	case projectSyncCmd.FullCommand():
		p := cli.Project{
			ID:       *projectSyncID,
			DomainID: *projectSyncDomain,
		}
		cli.RunSyncTask(&p)
	}
}

// custom parser magiggy for quota values
type quotas api.ServiceQuotas

func (q *quotas) Set(str string) error {
	rawQuotaVals := strings.Fields(str)

	for _, rq := range rawQuotaVals {
		// tmp holds the different components of a single parsed quota value. This makes it easier to refer
		// to individual components and pass them to the api.ServiceQuotas type
		var tmp struct {
			Service  string
			Resource string
			Value    uint64
			Unit     limes.Unit
		}

		// separate the value from the identifier
		idfVal := strings.SplitN(rq, "=", 2)
		if len(idfVal) != 2 {
			return fmt.Errorf("expected a quota value in the format: service/resource=value(unit), got '%s'", rq)
		}

		// separate the service from the resource
		srvRes := strings.SplitN(idfVal[0], "/", 2)
		if len(srvRes) != 2 {
			return fmt.Errorf("expected service/resource, got '%s'", idfVal[0])
		}
		tmp.Service = srvRes[0]
		tmp.Resource = srvRes[1]

		// separate the quota value from the unit and determine the type of unit, if one was given
		valStr := idfVal[1]
		tmp.Unit = limes.UnitNone
		// tmpVal allows easier string slicing
		tmpVal := idfVal[1]
		if len(tmpVal) > 1 && tmpVal[len(tmpVal)-1:] == "B" {
			tmp.Unit = limes.UnitBytes
			valStr = (tmpVal)[:(len(tmpVal) - 1)]
		}
		if len(tmpVal) > 3 {
			units := map[string]limes.Unit{
				"KiB": limes.UnitKibibytes,
				"MiB": limes.UnitMebibytes,
				"GiB": limes.UnitGibibytes,
				"TiB": limes.UnitTebibytes,
				"PiB": limes.UnitPebibytes,
				"EiB": limes.UnitExbibytes,
			}
			for unitStr, limesUnit := range units {
				if unitStr == tmpVal[len(tmpVal)-3:] {
					tmp.Unit = limesUnit
					valStr = (tmpVal)[:len(tmpVal)-3]
				}
			}
		}

		valUint, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			return fmt.Errorf("could not parse quota value: '%s'", valStr)
		}
		tmp.Value = valUint

		if _, exists := (*q)[tmp.Service]; !exists {
			(*q)[tmp.Service] = make(api.ResourceQuotas)
		}
		(*q)[tmp.Service][tmp.Resource] = limes.ValueWithUnit{tmp.Value, tmp.Unit}
	}

	return nil
}

func (q *quotas) String() string {
	return ""
}

func (q *quotas) IsCumulative() bool {
	return true
}

func Quota(s kingpin.Settings) (target *api.ServiceQuotas) {
	target = &api.ServiceQuotas{}
	s.SetValue((*quotas)(target))
	return
}

// magiggy continued...
type capacities map[string][]api.ResourceCapacity

func (c *capacities) Set(str string) error {
	rawCapacityVals := strings.Fields(str)

	for _, rc := range rawCapacityVals {
		// resCap holds the different components of a single parsed capacity value. This makes it easier to refer
		// to individual components and pass them to the map[string][]api.ResourceCapacity type
		var resCap api.ResourceCapacity

		// separate the value from the identifier
		idfVal := strings.SplitN(rc, "=", 2)
		if len(idfVal) != 2 {
			return fmt.Errorf("expected a capacity value in the format: service/resource=value(unit):comment, got '%s'", rc)
		}

		// separate the service from the resource
		srvRes := strings.SplitN(idfVal[0], "/", 2)
		if len(srvRes) != 2 {
			return fmt.Errorf("expected service/resource, got '%s'", idfVal[0])
		}
		tmpSrv := srvRes[0]
		resCap.Name = srvRes[1]

		//separate capacity value from comment
		valCom := strings.SplitN(idfVal[1], ":", 2)
		if len(valCom) != 2 {
			return fmt.Errorf("expected value(unit):comment, got '%s'", idfVal[1])
		}
		resCap.Comment = valCom[1]

		// separate the capacity value from the unit and determine the type of unit, if one was given
		valStr := valCom[0]
		tmpUnit := limes.UnitNone
		// tmpVal allows easier string slicing
		tmpVal := valCom[0]
		if len(tmpVal) > 1 && tmpVal[len(tmpVal)-1:] == "B" {
			tmpUnit = limes.UnitBytes
			valStr = (tmpVal)[:(len(tmpVal) - 1)]
		}
		if len(tmpVal) > 3 {
			units := map[string]limes.Unit{
				"KiB": limes.UnitKibibytes,
				"MiB": limes.UnitMebibytes,
				"GiB": limes.UnitGibibytes,
				"TiB": limes.UnitTebibytes,
				"PiB": limes.UnitPebibytes,
				"EiB": limes.UnitExbibytes,
			}
			for unitStr, limesUnit := range units {
				if unitStr == tmpVal[len(tmpVal)-3:] {
					tmpUnit = limesUnit
					valStr = (tmpVal)[:len(tmpVal)-3]
				}
			}
		}

		valInt, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return fmt.Errorf("could not parse quota value: '%s'", valStr)
		}
		resCap.Capacity = valInt
		resCap.Unit = &tmpUnit

		if _, exists := (*c)[tmpSrv]; !exists {
			(*c)[tmpSrv] = make([]api.ResourceCapacity, 0)
		}
		(*c)[tmpSrv] = append((*c)[tmpSrv], resCap)
	}

	return nil
}

func (c *capacities) String() string {
	return ""
}

func (c *capacities) IsCumulative() bool {
	return true
}

func Capacity(s kingpin.Settings) (target *map[string][]api.ResourceCapacity) {
	target = &map[string][]api.ResourceCapacity{}
	s.SetValue((*capacities)(target))
	return
}
