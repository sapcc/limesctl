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
	domainCluster = domainCmd.Flag("cluster", "Cluster ID. When this option is given, the domain must be identified by ID. Specifiying a domain name will not work.").Short('c').String()

	projectCmd     = app.Command("project", "Do some action on project(s).")
	projectCluster = projectCmd.Flag("cluster", "Cluster ID. When this option is given, the domain/project must be identified by ID. Specifiying a domain/project name will not work.").Short('c').String()

	osAuthURL           = app.Flag("os-auth-url", "Authentication URL.").PlaceHolder("OS_AUTH_URL").String()
	osUsername          = app.Flag("os-username", "Username").PlaceHolder("OS_USERNAME").String()
	osPassword          = app.Flag("os-password", "User's Password").PlaceHolder("OS_PASSWORD").String()
	osUserDomainID      = app.Flag("os-user-domain-id", "User's domain ID.").PlaceHolder("OS_USER_DOMAIN_ID").String()
	osUserDomainName    = app.Flag("os-user-domain-name", "User's domain name.").PlaceHolder("OS_USER_DOMAIN_NAME").String()
	osProjectID         = app.Flag("os-project-id", "Project ID to scope to.").PlaceHolder("OS_PROJECT_ID").String()
	osProjectName       = app.Flag("os-project-name", "Project name to scope to.").PlaceHolder("OS_PROJECT_NAME").String()
	osProjectDomainID   = app.Flag("os-project-domain-ID", "Domain ID containing project to scope to.").PlaceHolder("OS_PROJECT_DOMAIN_ID").String()
	osProjectDomainName = app.Flag("os-project-domain-name", "Domain name containing project to scope to.").PlaceHolder("OS_PROJECT_DOMAIN_NAME").String()

	area              = app.Flag("area", "Resource area.").String()
	service           = app.Flag("service", "Service type.").String()
	resource          = app.Flag("resource", "Resource name.").String()
	namesOutput       = app.Flag("names", "Show output with names instead of UUIDs.").Bool()
	longOutput        = app.Flag("long", "Show detailed output.").Bool()
	humanReadableVals = app.Flag("human-readable", "Show detailed output.").Bool()
	outputFmt         = app.Flag("format", "Output format (table, json, csv).").PlaceHolder("table").Short('f').Enum("table", "json", "csv")

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

	// parse all command-line args and flags
	cmdString := kingpin.MustParse(app.Parse(os.Args[1:]))

	// overwrite OpenStack variables
	setEnvUnlessEmpty("OS_AUTH_URL", *osAuthURL)
	setEnvUnlessEmpty("OS_USERNAME", *osUsername)
	setEnvUnlessEmpty("OS_PASSWORD", *osPassword)
	setEnvUnlessEmpty("OS_USER_DOMAIN_ID", *osUserDomainID)
	setEnvUnlessEmpty("OS_USER_DOMAIN_NAME", *osUserDomainName)
	setEnvUnlessEmpty("OS_PROJECT_ID", *osProjectID)
	setEnvUnlessEmpty("OS_PROJECT_NAME", *osProjectName)
	setEnvUnlessEmpty("OS_PROJECT_DOMAIN_ID", *osProjectDomainID)
	setEnvUnlessEmpty("OS_PROJECT_DOMAIN_NAME", *osProjectDomainName)

	// output and filter are initialized in advance with values that were provided
	// at the command-line. Later, we pass only the specific information that
	// is required by the operation
	filter := cli.Filter{
		Area:     *area,
		Service:  *service,
		Resource: *resource,
	}
	output := cli.Output{
		Names:         *namesOutput,
		Long:          *longOutput,
		HumanReadable: *humanReadableVals,
	}

	switch cmdString {
	case clusterListCmd.FullCommand():
		c := &cli.Cluster{
			Filter: filter,
			Output: output,
		}
		cli.RunListTask(c, *outputFmt)

	case clusterShowCmd.FullCommand():
		c := &cli.Cluster{
			ID:     *clusterShowID,
			Filter: filter,
			Output: output,
		}
		cli.RunGetTask(c, *outputFmt)

	case clusterSetCmd.FullCommand():
		// this manual check is required due to the order of the Args.
		// If the ID is not provided then the capacities get interpreted
		// as the ID and the error shown is not relevant to the context
		if strings.Contains(*clusterSetID, "=") {
			kingpin.Fatalf("required argument 'cluster-id' not provided, try --help")
		}
		c := &cli.Cluster{ID: *clusterSetID}
		cli.RunSetTask(c, clusterSetCaps)

	case domainListCmd.FullCommand():
		d := &cli.Domain{
			Filter: filter,
			Output: output,
		}
		d.Filter.Cluster = *domainCluster
		cli.RunListTask(d, *outputFmt)

	case domainShowCmd.FullCommand():
		// since gophercloud does not allow domain listing across
		// different clusters therefore we skip FindDomain(), if a cluster
		// was provided at the command-line
		var d *cli.Domain
		if *domainCluster == "" {
			var err error
			d, err = cli.FindDomain(*domainShowID)
			fatalIfErr(err)
		} else {
			d = &cli.Domain{ID: *domainShowID}
		}

		d.Filter = filter
		d.Filter.Cluster = *domainCluster
		d.Output = output
		cli.RunGetTask(d, *outputFmt)

	case domainSetCmd.FullCommand():
		if strings.Contains(*domainSetID, "=") {
			kingpin.Fatalf("required argument 'domain-id' not provided, try --help")
		}
		var d *cli.Domain
		if *domainCluster == "" {
			var err error
			d, err = cli.FindDomain(*domainSetID)
			fatalIfErr(err)
		} else {
			d = &cli.Domain{ID: *domainSetID}
		}

		d.Filter.Cluster = *domainCluster
		cli.RunSetTask(d, domainSetQuotas)

	case projectListCmd.FullCommand():
		var d *cli.Domain
		var err error
		if *projectCluster == "" {
			d, err = cli.FindDomain(*projectListDomain)
		} else {
			d, err = cli.FindDomainInCluster(*projectListDomain, *projectCluster)
		}
		fatalIfErr(err)

		p := &cli.Project{
			DomainID:   d.ID,
			DomainName: d.Name,
			Filter:     filter,
			Output:     output,
		}
		p.Filter.Cluster = *projectCluster
		cli.RunListTask(p, *outputFmt)

	case projectShowCmd.FullCommand():
		var p *cli.Project
		var err error
		if *projectCluster == "" {
			p, err = cli.FindProject(*projectShowID, *projectShowDomain)
		} else {
			p, err = cli.FindProjectInCluster(*projectShowID, *projectShowDomain, *projectCluster)
		}
		fatalIfErr(err)

		p.Filter = filter
		p.Filter.Cluster = *projectCluster
		p.Output = output
		cli.RunGetTask(p, *outputFmt)

	case projectSetCmd.FullCommand():
		if strings.Contains(*projectSetID, "=") {
			kingpin.Fatalf("required argument 'project-id' not provided, try --help")
		}
		var p *cli.Project
		var err error
		if *projectCluster == "" {
			p, err = cli.FindProject(*projectSetID, *projectSetDomain)
		} else {
			p, err = cli.FindProjectInCluster(*projectSetID, *projectSetDomain, *projectCluster)
		}
		fatalIfErr(err)

		p.Filter.Cluster = *projectCluster
		cli.RunSetTask(p, projectSetQuotas)

	case projectSyncCmd.FullCommand():
		var p *cli.Project
		var err error
		if *projectCluster == "" {
			p, err = cli.FindProject(*projectSyncID, *projectSyncDomain)
		} else {
			p, err = cli.FindProjectInCluster(*projectSyncID, *projectSyncDomain, *projectCluster)
		}
		fatalIfErr(err)

		p.Filter.Cluster = *projectCluster
		cli.RunSyncTask(p)
	}
}

func setEnvUnlessEmpty(env, val string) {
	if val == "" {
		return
	}

	os.Setenv(env, val)
}

func fatalIfErr(err error) {
	if err == nil {
		return
	}

	kingpin.Fatalf(err.Error())
}
