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
	"github.com/sapcc/limesctl/internal/core"
	"github.com/sapcc/limesctl/internal/errors"
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
	osProjectDomainID   = app.Flag("os-project-domain-id", "Domain ID containing project to scope to.").PlaceHolder("OS_PROJECT_DOMAIN_ID").String()
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
	clusterSetCaps = QuotaList(clusterSetCmd.Arg("capacities", "Capacities to change. Format: service/resource=value(unit):\"comment\"").Required())

	domainListCmd = domainCmd.Command("list", "Query data for all the domains. Requires a cloud-admin token.")

	domainShowCmd = domainCmd.Command("show", "Query data for a specific domain. Requires a domain-admin token.")
	domainShowID  = domainShowCmd.Arg("domain-id", "Domain ID (name/UUID).").Required().String()

	domainSetCmd    = domainCmd.Command("set", "Change resource(s) quota for a domain. Requires a cloud-admin token.")
	domainSetID     = domainSetCmd.Arg("domain-id", "Domain ID (name/UUID).").Required().String()
	domainSetQuotas = QuotaList(domainSetCmd.Arg("quotas", "Quotas to change. Format: service/resource=value(unit)").Required())

	projectListCmd    = projectCmd.Command("list", "Query data for all the projects in a domain. Requires a domain-admin token.")
	projectListDomain = projectListCmd.Flag("domain", "Domain ID.").Short('d').Required().String()

	projectShowCmd    = projectCmd.Command("show", "Query data for a specific project in a domain. Requires project member permissions.")
	projectShowDomain = projectShowCmd.Flag("domain", "Domain ID.").Short('d').String()
	projectShowID     = projectShowCmd.Arg("project-id", "Project ID (name/UUID).").Required().String()

	projectSetCmd    = projectCmd.Command("set", "Change resource(s) quota for a project. Requires a domain-admin token.")
	projectSetDomain = projectSetCmd.Flag("domain", "Domain ID.").Short('d').String()
	projectSetID     = projectSetCmd.Arg("project-id", "Project ID (name/UUID).").Required().String()
	projectSetQuotas = QuotaList(projectSetCmd.Arg("quotas", "Quotas to change. Format: service/resource=value(unit)").Required())

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
	filter := core.Filter{
		Area:     *area,
		Service:  *service,
		Resource: *resource,
	}
	output := core.Output{
		Names:         *namesOutput,
		Long:          *longOutput,
		HumanReadable: *humanReadableVals,
	}

	switch cmdString {
	case clusterListCmd.FullCommand():
		c := &core.Cluster{
			Filter: filter,
			Output: output,
		}
		core.RunListTask(c, *outputFmt)

	case clusterShowCmd.FullCommand():
		c := &core.Cluster{
			ID:     *clusterShowID,
			Filter: filter,
			Output: output,
		}
		core.RunGetTask(c, *outputFmt)

	case clusterSetCmd.FullCommand():
		// this manual check is required due to the order of the Args.
		// If the ID is not provided then the capacities get interpreted
		// as the ID and the error shown is not relevant to the context
		if strings.Contains(*clusterSetID, "=") {
			errors.Handle(errors.New("required argument 'cluster-id' not provided, try --help"))
		}

		c := &core.Cluster{ID: *clusterSetID}
		q, err := core.ParseRawQuotas(c, clusterSetCaps, false)
		errors.Handle(err)
		core.RunSetTask(c, q)

	case domainListCmd.FullCommand():
		d := &core.Domain{
			Filter: filter,
			Output: output,
		}
		d.Filter.Cluster = *domainCluster
		core.RunListTask(d, *outputFmt)

	case domainShowCmd.FullCommand():
		d, err := core.FindDomain(*domainShowID, *domainCluster)
		errors.Handle(err)

		d.Filter = filter
		d.Filter.Cluster = *domainCluster
		d.Output = output
		core.RunGetTask(d, *outputFmt)

	case domainSetCmd.FullCommand():
		if strings.Contains(*domainSetID, "=") {
			errors.Handle(errors.New("required argument 'domain-id' not provided, try --help"))
		}
		d, err := core.FindDomain(*domainSetID, *domainCluster)
		errors.Handle(err)

		d.Filter.Cluster = *domainCluster
		q, err := core.ParseRawQuotas(d, domainSetQuotas, false)
		errors.Handle(err)
		core.RunSetTask(d, q)

	case projectListCmd.FullCommand():
		d, err := core.FindDomain(*projectListDomain, *projectCluster)
		errors.Handle(err)

		p := &core.Project{
			DomainID:   d.ID,
			DomainName: d.Name,
			Filter:     filter,
			Output:     output,
		}
		p.Filter.Cluster = *projectCluster
		core.RunListTask(p, *outputFmt)

	case projectShowCmd.FullCommand():
		var p *core.Project
		var err error
		if *projectCluster == "" {
			p, err = core.FindProject(*projectShowID, *projectShowDomain)
		} else {
			p, err = core.FindProjectInCluster(*projectShowID, *projectShowDomain, *projectCluster)
		}
		errors.Handle(err)

		p.Filter = filter
		p.Filter.Cluster = *projectCluster
		p.Output = output
		core.RunGetTask(p, *outputFmt)

	case projectSetCmd.FullCommand():
		if strings.Contains(*projectSetID, "=") {
			errors.Handle(errors.New("required argument 'project-id' not provided, try --help"))
		}
		var p *core.Project
		var err error
		if *projectCluster == "" {
			p, err = core.FindProject(*projectSetID, *projectSetDomain)
		} else {
			p, err = core.FindProjectInCluster(*projectSetID, *projectSetDomain, *projectCluster)
		}
		errors.Handle(err)

		p.Filter.Cluster = *projectCluster
		q, err := core.ParseRawQuotas(p, projectSetQuotas, false)
		errors.Handle(err)
		core.RunSetTask(p, q)

	case projectSyncCmd.FullCommand():
		var p *core.Project
		var err error
		if *projectCluster == "" {
			p, err = core.FindProject(*projectSyncID, *projectSyncDomain)
		} else {
			p, err = core.FindProjectInCluster(*projectSyncID, *projectSyncDomain, *projectCluster)
		}
		errors.Handle(err)

		p.Filter.Cluster = *projectCluster
		core.RunSyncTask(p)
	}
}

func setEnvUnlessEmpty(env, val string) {
	if val == "" {
		return
	}

	os.Setenv(env, val)
}

type rawQuotas core.RawQuotas

// Set implements the kingpin.Value interface.
func (rq *rawQuotas) Set(value string) error {
	*rq = append(*rq, strings.TrimSpace(value))
	return nil
}

// String implements the kingpin.Value interface.
func (rq *rawQuotas) String() string {
	return ""
}

// IsCumulative allows consumption of remaining command line arguments.
func (rq *rawQuotas) IsCumulative() bool {
	return true
}

// QuotaList appends the raw quota values given at the command line to the
// aggregate core.RawQuotas list.
func QuotaList(s kingpin.Settings) (target *core.RawQuotas) {
	target = new(core.RawQuotas)
	s.SetValue((*rawQuotas)(target))
	return
}
