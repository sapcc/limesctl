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
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin"
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
	clusterSetCaps = ParseQuotas(clusterSetCmd.Arg("quota", "Quota value(s) to change. Format: service/resource=value(unit)").Required())

	domainListCmd = domainCmd.Command("list", "Query data for all the domains. Requires a cloud-admin token.")

	domainShowCmd = domainCmd.Command("show", "Query data for a specific domain. Requires a domain-admin token.")
	domainShowID  = domainShowCmd.Arg("domain-id", "Domain ID (name/UUID).").Required().String()

	domainSetCmd    = domainCmd.Command("set", "Change resource(s) quota for a domain. Requires a cloud-admin token.")
	domainSetID     = domainSetCmd.Arg("domain-id", "Domain ID (name/UUID).").Required().String()
	domainSetQuotas = ParseQuotas(domainSetCmd.Arg("quota", "Quota value(s) to change. Format: service/resource=value(unit)").Required())

	projectListCmd    = projectCmd.Command("list", "Query data for all the projects in a domain. Requires a domain-admin token.")
	projectListDomain = projectListCmd.Flag("domain", "Domain ID.").Short('d').Required().String()

	projectShowCmd    = projectCmd.Command("show", "Query data for a specific project in a domain. Requires project member permissions.")
	projectShowDomain = projectShowCmd.Flag("domain", "Domain ID.").Short('d').String()
	projectShowID     = projectShowCmd.Arg("project-id", "Project ID (name/UUID).").Required().String()

	projectSetCmd    = projectCmd.Command("set", "Change resource(s) quota for a project. Requires a domain-admin token.")
	projectSetDomain = projectSetCmd.Flag("domain", "Domain ID.").Short('d').String()
	projectSetID     = projectSetCmd.Arg("project-id", "Project ID (name/UUID).").Required().String()
	projectSetQuotas = ParseQuotas(projectSetCmd.Arg("quota", "Quota value(s) to change. Format: service/resource=value(unit)").Required())

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
		}
		cli.RunSetTask(&c, clusterSetCaps)

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
		}
		cli.RunSetTask(&d, domainSetQuotas)

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
		}
		cli.RunSetTask(&p, projectSetQuotas)
	case projectSyncCmd.FullCommand():
		p := cli.Project{
			ID:       *projectSyncID,
			DomainID: *projectSyncDomain,
		}
		cli.RunSyncTask(&p)
	}
}

type quotas cli.Quotas

// Set implements the kingpin.Value interface
func (q *quotas) Set(value string) error {
	value = strings.TrimSpace(value)
	// tmp holds the different components of a single parsed quota value. This makes it easier to refer
	// to individual components and pass them to the cli.Quotas map
	var tmp cli.Resource

	// separate the quota value from its identifier
	idfVal := strings.SplitN(value, "=", 2)
	if len(idfVal) != 2 {
		return fmt.Errorf("expected a quota in the format: service/resource=value(unit), got '%s'", value)
	}

	// separate service and resource
	srvRes := strings.SplitN(idfVal[0], "/", 2)
	if len(srvRes) != 2 {
		return fmt.Errorf("expected service/resource, got '%s'", idfVal[0])
	}
	srv := srvRes[0]
	tmp.Name = srvRes[1]

	// separate quota value and comment (if one was given)
	valCom := strings.SplitN(idfVal[1], ":", 2)
	if len(valCom) > 1 {
		if valCom[1] != "" {
			tmp.Comment = strings.TrimSpace(valCom[1])
		}
	}

	// separate quota's value from its unit (if one was given)
	rx := regexp.MustCompile(`^([0-9]+)([A-Za-z]+)?$`)
	match := rx.MatchString(valCom[0])
	if !match {
		return fmt.Errorf("expected a quota value with optional unit in the format: 123Unit, got '%s'", valCom[0])
	}

	// rxMatchedList: []string{"entire regex matched string", "quota value", "unit (empty, if no unit given)"}
	rxMatchedList := rx.FindStringSubmatch(valCom[0])
	intVal, err := strconv.ParseInt(rxMatchedList[1], 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse quota value: '%s'", rxMatchedList[1])
	}
	tmp.Value = intVal

	tmp.Unit = limes.UnitNone
	if rxMatchedList[2] != "" {
		switch rxMatchedList[2] {
		case "B":
			tmp.Unit = limes.UnitBytes
		case "KiB":
			tmp.Unit = limes.UnitKibibytes
		case "MiB":
			tmp.Unit = limes.UnitMebibytes
		case "GiB":
			tmp.Unit = limes.UnitGibibytes
		case "TiB":
			tmp.Unit = limes.UnitTebibytes
		case "PiB":
			tmp.Unit = limes.UnitPebibytes
		case "EiB":
			tmp.Unit = limes.UnitExbibytes
		default:
			return fmt.Errorf("acceptable units: ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB'], got '%s'", rxMatchedList[2])
		}
	}

	if _, exists := (*q)[srv]; !exists {
		(*q)[srv] = make([]cli.Resource, 0)
	}
	(*q)[srv] = append((*q)[srv], tmp)

	return nil
}

// String implements the kingpin.Value interface
func (q *quotas) String() string {
	return ""
}

// IsCumulative allows consumption of remaining command line arguments.
func (q *quotas) IsCumulative() bool {
	return true
}

// ParseQuotas parses a command line argument to a quota value and assigns it to the
// aggregate cli.Quotas map.
func ParseQuotas(s kingpin.Settings) (target *cli.Quotas) {
	target = &cli.Quotas{}
	s.SetValue((*quotas)(target))
	return
}
