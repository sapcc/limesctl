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
	"os"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/sapcc/limesctl/pkg/cli"
)

var (
	// defined by the Makefile at compile time
	version string

	app = kingpin.New("limesctl", "CLI client for Limes.")
	// first-level commands and flags
	clusterCmd = app.Command("cluster", "Do some action on cluster(s).")

	domainCmd     = app.Command("domain", "Do some action on domain(s).")
	domainCluster = domainCmd.Flag("cluster", "Cluster ID.").Short('c').String()

	projectCmd     = app.Command("project", "Do some action on project(s).")
	projectCluster = projectCmd.Flag("cluster", "Cluster ID.").Short('c').String()

	area        = app.Flag("area", "Resource area.").String()
	service     = app.Flag("service", "Service type.").String()
	resource    = app.Flag("resource", "Resource name.").String()
	namesOutput = app.Flag("names", "Show output with names instead of UUIDs.").Bool()
	longOutput  = app.Flag("long", "Show detailed output.").Bool()
	outputFmt   = app.Flag("format", "Output format (table, json, csv).").PlaceHolder("table").Short('f').Enum("table", "json", "csv")

	// second-level subcommands and their flags/args
	clusterListCmd = clusterCmd.Command("list", "Query data for all the clusters. Requires a cloud-admin token.")

	clusterShowCmd = clusterCmd.Command("show", "Query data for a specific cluster. Use 'current' to show information regarding the current cluster. Requires a cloud-admin token.")
	clusterShowID  = clusterShowCmd.Arg("cluster-id", "Cluster ID.").Required().String()

	clusterSetCmd  = clusterCmd.Command("set", "Change resource(s) quota for a cluster. Use 'current' to show information regarding the current cluster. Requires a cloud-admin token.")
	clusterSetID   = clusterSetCmd.Arg("cluster-id", "Cluster ID.").Required().String()
	clusterSetCaps = cli.ParseQuotas(clusterSetCmd.Arg("capacities", "Capacities to change. Format: service/resource=value(unit):\"comment\"").Required())

	domainListCmd = domainCmd.Command("list", "Query data for all the domains. Requires a cloud-admin token.")

	domainShowCmd = domainCmd.Command("show", "Query data for a specific domain. Requires a domain-admin token.")
	domainShowID  = domainShowCmd.Arg("domain-id", "Domain ID (name/UUID).").Required().String()

	domainSetCmd    = domainCmd.Command("set", "Change resource(s) quota for a domain. Requires a cloud-admin token.")
	domainSetID     = domainSetCmd.Arg("domain-id", "Domain ID (name/UUID).").Required().String()
	domainSetQuotas = cli.ParseQuotas(domainSetCmd.Arg("quotas", "Quotas to change. Format: service/resource=value(unit)").Required())

	projectListCmd    = projectCmd.Command("list", "Query data for all the projects in a domain. Requires a domain-admin token.")
	projectListDomain = projectListCmd.Flag("domain", "Domain ID.").Short('d').Required().String()

	projectShowCmd    = projectCmd.Command("show", "Query data for a specific project in a domain. Requires project member permissions.")
	projectShowDomain = projectShowCmd.Flag("domain", "Domain ID.").Short('d').String()
	projectShowID     = projectShowCmd.Arg("project-id", "Project ID (name/UUID).").Required().String()

	projectSetCmd    = projectCmd.Command("set", "Change resource(s) quota for a project. Requires a domain-admin token.")
	projectSetDomain = projectSetCmd.Flag("domain", "Domain ID.").Short('d').String()
	projectSetID     = projectSetCmd.Arg("project-id", "Project ID (name/UUID).").Required().String()
	projectSetQuotas = cli.ParseQuotas(projectSetCmd.Arg("quotas", "Quotas to change. Format: service/resource=value(unit)").Required())

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
		c := &cli.Cluster{
			Opts: cli.Options{
				Long:     *longOutput,
				Area:     *area,
				Service:  *service,
				Resource: *resource,
			},
		}
		cli.RunListTask(c, *outputFmt)
	case clusterShowCmd.FullCommand():
		c := &cli.Cluster{
			ID: *clusterShowID,
			Opts: cli.Options{
				Long:     *longOutput,
				Area:     *area,
				Service:  *service,
				Resource: *resource,
			},
		}
		cli.RunGetTask(c, *outputFmt)
	case clusterSetCmd.FullCommand():
		// this manual check is required due to the order of the Args.
		// If the ID is not provided then the capacities get interpreted
		// as the ID and the error shown is not relevant to the context
		if strings.Contains(*clusterSetID, "=") {
			kingpin.Fatalf("required argument 'cluster-id' not provided, try --help")
		}
		c := &cli.Cluster{
			ID: *clusterSetID,
		}
		cli.RunSetTask(c, clusterSetCaps)
	case domainListCmd.FullCommand():
		d := &cli.Domain{
			Opts: cli.Options{
				Names:    *namesOutput,
				Long:     *longOutput,
				Cluster:  *domainCluster,
				Area:     *area,
				Service:  *service,
				Resource: *resource,
			},
		}
		cli.RunListTask(d, *outputFmt)
	case domainShowCmd.FullCommand():
		d, err := cli.FindDomain(*domainShowID)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		d.Opts = cli.Options{
			Names:    *namesOutput,
			Long:     *longOutput,
			Cluster:  *domainCluster,
			Area:     *area,
			Service:  *service,
			Resource: *resource,
		}
		cli.RunGetTask(d, *outputFmt)
	case domainSetCmd.FullCommand():
		if strings.Contains(*domainSetID, "=") {
			kingpin.Fatalf("required argument 'domain-id' not provided, try --help")
		}
		d, err := cli.FindDomain(*domainSetID)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		d.Opts.Cluster = *domainCluster
		cli.RunSetTask(d, domainSetQuotas)
	case projectListCmd.FullCommand():
		d, err := cli.FindDomain(*projectListDomain)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		p := &cli.Project{
			DomainID:   d.ID,
			DomainName: d.Name,
			Opts: cli.Options{
				Names:    *namesOutput,
				Long:     *longOutput,
				Cluster:  *projectCluster,
				Area:     *area,
				Service:  *service,
				Resource: *resource,
			},
		}
		cli.RunListTask(p, *outputFmt)
	case projectShowCmd.FullCommand():
		p, err := cli.FindProject(*projectShowID, *projectShowDomain)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		p.Opts = cli.Options{
			Names:    *namesOutput,
			Long:     *longOutput,
			Cluster:  *projectCluster,
			Area:     *area,
			Service:  *service,
			Resource: *resource,
		}
		cli.RunGetTask(p, *outputFmt)
	case projectSetCmd.FullCommand():
		if strings.Contains(*projectSetID, "=") {
			kingpin.Fatalf("required argument 'project-id' not provided, try --help")
		}
		p, err := cli.FindProject(*projectSetID, *projectSetDomain)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		p.Opts.Cluster = *projectCluster
		cli.RunSetTask(p, projectSetQuotas)
	case projectSyncCmd.FullCommand():
		p, err := cli.FindProject(*projectSyncID, *projectSyncDomain)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		p.Opts.Cluster = *projectCluster
		cli.RunSyncTask(p)
	}
}
