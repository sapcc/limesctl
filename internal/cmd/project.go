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

package cmd

import (
	"fmt"

	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/go-api-declarations/limes"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
	ratesProjects "github.com/sapcc/gophercloud-sapcc/rates/v1/projects"
	"github.com/sapcc/gophercloud-sapcc/resources/v1/projects"
	"github.com/spf13/cobra"

	"github.com/sapcc/limesctl/v3/internal/auth"
	"github.com/sapcc/limesctl/v3/internal/core"
	"github.com/sapcc/limesctl/v3/internal/util"
)

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Do some action at project level",
		Args:  cobra.NoArgs,
	}
	// Flags
	doNotSortFlags(cmd)
	// Subcommands
	cmd.AddCommand(newProjectListCmd().Command)
	cmd.AddCommand(newProjectListRatesCmd().Command)
	cmd.AddCommand(newProjectShowCmd().Command)
	cmd.AddCommand(newProjectShowRatesCmd().Command)
	cmd.AddCommand(newProjectSetCmd().Command)
	cmd.AddCommand(newProjectSyncCmd().Command)
	return cmd
}

type projectFlags struct {
	DomainNameOrID string `short:"d" name:"domain" help:"Name or ID of the domain."`
}

func (pf *projectFlags) AddToCmd(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&pf.DomainNameOrID, "domain", "d", "", "name or ID of the domain")
}

func (pf *projectFlags) validateWithNameID(nameOrID string) error {
	if pf.DomainNameOrID != "" && nameOrID == "" {
		return errors.New("project name or ID is required when using the '--domain' flag")
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Project list.

type projectListCmd struct {
	*cobra.Command

	projectFlags
	resourceFilterFlags
	resourceOutputFmtFlags
}

func newProjectListCmd() *projectListCmd {
	projectList := &projectListCmd{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Display resource usage data for all the projects in a domain",
		Long: `Display resource usage data for all the projects in a domain.

The project name/ID is optional by default and limesctl will get the project
from current scope. However, if '--domain' flag is used then either project
name or ID is required.

This command requires a domain-admin token.`,
		Args:    cobra.NoArgs,
		PreRunE: authWithLimesResources,
		RunE:    projectList.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	projectList.projectFlags.AddToCmd(cmd)
	projectList.resourceFilterFlags.AddToCmd(cmd)
	projectList.resourceOutputFmtFlags.AddToCmd(cmd)

	projectList.Command = cmd
	return projectList
}

func (p *projectListCmd) Run(_ *cobra.Command, _ []string) error {
	outputOpts, err := p.resourceOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	domainName := ""
	domainID, err := auth.FindDomainID(identityClient, p.DomainNameOrID)
	if err == nil {
		domainName, err = auth.FindDomainName(identityClient, domainID)
	}
	if err != nil {
		return err
	}

	res := projects.List(limesResourcesClient, domainID, projects.ListOpts{
		Areas:     p.areas,
		Services:  p.services,
		Resources: p.resources,
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get project reports")
	}

	if p.format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesReps, err := res.ExtractProjects()
	if err != nil {
		return util.WrapError(err, "could not extract project reports")
	}

	return writeReports(outputOpts,
		core.LimesProjectResourcesToReportRenderer(limesReps, domainID, domainName, false)...)
}

///////////////////////////////////////////////////////////////////////////////
// Project list rates.

type projectListRatesCmd struct {
	*cobra.Command

	projectFlags
	rateFilterFlags
	rateOutputFmtFlags
}

func newProjectListRatesCmd() *projectListRatesCmd {
	projectListRates := &projectListRatesCmd{}
	cmd := &cobra.Command{
		Use:   "list-rates",
		Short: "Display rate limits for all the projects in a domain",
		Long: `Display rate limits for all the projects in a domain.

The project name/ID is optional by default and limesctl will get the project
from current scope. However, if '--domain' flag is used then either project
name or ID is required.

This command requires a domain-admin token.`,
		Args:    cobra.NoArgs,
		PreRunE: authWithLimesRates,
		RunE:    projectListRates.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	projectListRates.projectFlags.AddToCmd(cmd)
	projectListRates.rateFilterFlags.AddToCmd(cmd)
	projectListRates.rateOutputFmtFlags.AddToCmd(cmd)

	projectListRates.Command = cmd
	return projectListRates
}

func (p *projectListRatesCmd) Run(_ *cobra.Command, _ []string) error {
	outputOpts, err := p.rateOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	domainName := ""
	domainID, err := auth.FindDomainID(identityClient, p.DomainNameOrID)
	if err == nil {
		domainName, err = auth.FindDomainName(identityClient, domainID)
	}
	if err != nil {
		return err
	}

	res := ratesProjects.List(limesRatesClient, domainID, ratesProjects.ReadOpts{
		Services: p.services,
		Areas:    p.areas,
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get project reports")
	}

	if p.format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesReps, err := res.ExtractProjects()
	if err != nil {
		return util.WrapError(err, "could not extract project reports")
	}

	return writeReports(outputOpts,
		core.LimesProjectRatesToReportRenderer(limesReps, domainID, domainName, true)...)
}

///////////////////////////////////////////////////////////////////////////////
// Project show.

type projectShowCmd struct {
	*cobra.Command

	projectFlags
	resourceFilterFlags
	resourceOutputFmtFlags
}

func newProjectShowCmd() *projectShowCmd {
	projectShow := &projectShowCmd{}
	cmd := &cobra.Command{
		Use:   "show [name or ID]",
		Short: "Display resource usage data for a specific project",
		Long: `Display resource usage data for a specific project

The project name/ID is optional by default and limesctl will get the project
from current scope. However, if '--domain' flag is used then either project
name or ID is required.

This command requires a project member permissions.`,
		Args:    cobra.MaximumNArgs(1),
		PreRunE: authWithLimesResources,
		RunE:    projectShow.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	projectShow.projectFlags.AddToCmd(cmd)
	projectShow.resourceFilterFlags.AddToCmd(cmd)
	projectShow.resourceOutputFmtFlags.AddToCmd(cmd)

	projectShow.Command = cmd
	return projectShow
}

func (p *projectShowCmd) Run(_ *cobra.Command, args []string) error {
	nameOrID := ""
	if len(args) > 0 {
		nameOrID = args[0]
	}
	err := p.projectFlags.validateWithNameID(nameOrID)
	if err != nil {
		return err
	}

	outputOpts, err := p.resourceOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	pInfo, err := auth.FindProject(identityClient, p.DomainNameOrID, nameOrID)
	if err != nil {
		return err
	}

	res := projects.Get(limesResourcesClient, pInfo.DomainID, pInfo.ID, projects.GetOpts{
		Areas:     p.areas,
		Services:  p.services,
		Resources: p.resources,
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get project report")
	}

	if p.format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesRep, err := res.Extract()
	if err != nil {
		return util.WrapError(err, "could not extract project report")
	}

	return writeReports(outputOpts, core.ProjectResourcesReport{
		ProjectReport: limesRep,
		DomainID:      pInfo.DomainID,
		DomainName:    pInfo.DomainName,
	})
}

///////////////////////////////////////////////////////////////////////////////
// Project show rates.

type projectShowRatesCmd struct {
	*cobra.Command

	projectFlags
	rateFilterFlags
	rateOutputFmtFlags
}

func newProjectShowRatesCmd() *projectShowRatesCmd {
	projectShowRates := &projectShowRatesCmd{}
	cmd := &cobra.Command{
		Use:   "show-rates [name or ID]",
		Short: "Display rate limits for a specific project",
		Long: `Display rate limits for a specific project.

The project name/ID is optional by default and limesctl will get the project
from current scope. However, if '--domain' flag is used then either project
name or ID is required.

This command requires a project member permissions.`,
		Args:    cobra.MaximumNArgs(1),
		PreRunE: authWithLimesRates,
		RunE:    projectShowRates.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	projectShowRates.projectFlags.AddToCmd(cmd)
	projectShowRates.rateFilterFlags.AddToCmd(cmd)
	projectShowRates.rateOutputFmtFlags.AddToCmd(cmd)

	projectShowRates.Command = cmd
	return projectShowRates
}

func (p *projectShowRatesCmd) Run(_ *cobra.Command, args []string) error {
	nameOrID := ""
	if len(args) > 0 {
		nameOrID = args[0]
	}
	err := p.projectFlags.validateWithNameID(nameOrID)
	if err != nil {
		return err
	}

	outputOpts, err := p.rateOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	pInfo, err := auth.FindProject(identityClient, p.DomainNameOrID, nameOrID)
	if err != nil {
		return err
	}

	res := ratesProjects.Get(limesRatesClient, pInfo.DomainID, pInfo.ID, ratesProjects.ReadOpts{
		Services: p.services,
		Areas:    p.areas,
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get project report")
	}

	if p.format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesRep, err := res.Extract()
	if err != nil {
		return util.WrapError(err, "could not extract project report")
	}

	return writeReports(outputOpts, core.ProjectRatesReport{
		ProjectReport: limesRep,
		DomainID:      pInfo.DomainID,
		DomainName:    pInfo.DomainName,
	})
}

///////////////////////////////////////////////////////////////////////////////
// Project set.

//nolint:lll
type projectSetCmd struct {
	*cobra.Command

	projectFlags
	quotas []string
}

func newProjectSetCmd() *projectSetCmd {
	projectSet := &projectSetCmd{}
	cmd := &cobra.Command{
		Use:   "set [name or ID]",
		Short: "Change resource quota values for a specific project",
		Long: `Change resource quota values for a specific project.

The project name/ID is optional by default and limesctl will get the project
from current scope. However, if '--domain' flag is used then either project
name or ID is required.

For relative quota adjustment, use one of the following operators: [+=, -=, *=, /=].

This command requires a domain-admin token.`,
		Example: makeExamplesString([]string{
			`limesctl project set --quotas="compute/cores=200,compute/ram=50GiB"`,
			`limesctl project set -q compute/cores=200 -q compute/ram=50GiB (you can give the flag multiple times and use flag shorthand '-q')`,
			`limesctl project set -q compute/cores*=2 -q compute/ram+=10GiB (relative quota update)`,
			`limesctl project set -q object-store/capacity=1TiB (you can also use a unit other than the service's default, e.g. object-store uses 'B' by default but we use 'TiB' here)`,
			`limesctl project set -q object-store/capacity-=0.25TiB (fractional values are also supported)`,
		}),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: authWithLimesResources,
		RunE:    projectSet.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	cmd.Flags().StringSliceVarP(&projectSet.quotas, "quotas", "q", nil, "new quota values (comma separated list)")
	cmd.MarkFlagRequired("quotas") //nolint: errcheck
	projectSet.projectFlags.AddToCmd(cmd)

	projectSet.Command = cmd
	return projectSet
}

func (p *projectSetCmd) Run(_ *cobra.Command, args []string) error {
	nameOrID := ""
	if len(args) > 0 {
		nameOrID = args[0]
	}
	err := p.projectFlags.validateWithNameID(nameOrID)
	if err != nil {
		return err
	}

	pInfo, err := auth.FindProject(identityClient, p.DomainNameOrID, nameOrID)
	if err != nil {
		return err
	}

	resQuotas, err := getProjectResourceQuotas(limesResourcesClient, pInfo)
	if err != nil {
		return util.WrapError(err, "could not get default units")
	}
	qc, err := parseToQuotaRequest(resQuotas, p.quotas)
	if err != nil {
		return util.WrapError(err, "could not parse quota values")
	}

	warn, err := projects.Update(limesResourcesClient, pInfo.DomainID, pInfo.ID, projects.UpdateOpts{
		Services: qc,
	}).Extract()
	if err != nil {
		return util.WrapError(err, "could not set new quotas for project")
	}

	if warn != nil {
		fmt.Println(string(warn))
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Project sync.

type projectSyncCmd struct {
	*cobra.Command

	projectFlags
}

func newProjectSyncCmd() *projectSyncCmd {
	projectSync := &projectSyncCmd{}
	cmd := &cobra.Command{
		Use:   "sync [name or ID]",
		Short: "Sync a specific project's resource data",
		Long: `Schedule a sync job that pulls quota and usage data for a specific project from
the backing services into Limes' local database.

The project name/ID is optional by default and limesctl will get the project
from current scope. However, if '--domain' flag is used then either project
name or ID is required.

This command requires a project-admin token.`,
		Args:    cobra.MaximumNArgs(1),
		PreRunE: authWithLimesResources,
		RunE:    projectSync.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	projectSync.projectFlags.AddToCmd(cmd)

	projectSync.Command = cmd
	return projectSync
}

func (p *projectSyncCmd) Run(_ *cobra.Command, args []string) error {
	nameOrID := ""
	if len(args) > 0 {
		nameOrID = args[0]
	}
	err := p.projectFlags.validateWithNameID(nameOrID)
	if err != nil {
		return err
	}

	pInfo, err := auth.FindProject(identityClient, p.DomainNameOrID, nameOrID)
	if err != nil {
		return err
	}

	err = projects.Sync(limesResourcesClient, pInfo.DomainID, pInfo.ID).ExtractErr()
	if err != nil {
		return util.WrapError(err, "could not sync project")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Helper functions.

func getProjectResourceQuotas(limesClient *gophercloud.ServiceClient, pInfo *auth.ProjectInfo) (resourceQuotas, error) {
	rep, err := projects.Get(limesClient, pInfo.DomainID, pInfo.ID, projects.GetOpts{}).Extract()
	if err != nil {
		return nil, err
	}

	result := make(resourceQuotas)
	for srv, srvReport := range rep.Services {
		for res, resReport := range srvReport.Resources {
			if _, ok := result[srv]; !ok {
				result[srv] = make(map[limesresources.ResourceName]limes.ValueWithUnit)
			}
			var val uint64
			if resReport.Quota != nil {
				val = *resReport.Quota
			}
			result[srv][res] = limes.ValueWithUnit{Value: val, Unit: resReport.ResourceInfo.Unit}
		}
	}
	return result, nil
}
