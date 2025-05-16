// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"

	"github.com/sapcc/go-api-declarations/limes"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
	ratesProjects "github.com/sapcc/gophercloud-sapcc/v2/rates/v1/projects"
	"github.com/sapcc/gophercloud-sapcc/v2/resources/v1/projects"
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

	projectFlags   projectFlags
	filterFlags    resourceFilterFlags
	outputFmtFlags resourceOutputFmtFlags
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
	projectList.filterFlags.AddToCmd(cmd)
	projectList.outputFmtFlags.AddToCmd(cmd)

	projectList.Command = cmd
	return projectList
}

func (p *projectListCmd) Run(cmd *cobra.Command, _ []string) error {
	outputOpts, err := p.outputFmtFlags.validate()
	if err != nil {
		return err
	}

	domainName := ""
	domainID, err := auth.FindDomainID(cmd.Context(), identityClient, p.projectFlags.DomainNameOrID)
	if err == nil {
		domainName, err = auth.FindDomainName(cmd.Context(), identityClient, domainID)
	}
	if err != nil {
		return err
	}

	res := projects.List(cmd.Context(), limesResourcesClient, domainID, projects.ListOpts{
		Areas:     p.filterFlags.areas,
		Services:  util.CastStringsTo[limes.ServiceType](p.filterFlags.services),
		Resources: util.CastStringsTo[limesresources.ResourceName](p.filterFlags.resources),
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get project reports")
	}

	if p.outputFmtFlags.format == core.OutputFormatJSON {
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

	projectFlags   projectFlags
	filterFlags    rateFilterFlags
	outputFmtFlags rateOutputFmtFlags
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
	projectListRates.filterFlags.AddToCmd(cmd)
	projectListRates.outputFmtFlags.AddToCmd(cmd)

	projectListRates.Command = cmd
	return projectListRates
}

func (p *projectListRatesCmd) Run(cmd *cobra.Command, _ []string) error {
	outputOpts, err := p.outputFmtFlags.validate()
	if err != nil {
		return err
	}

	domainName := ""
	domainID, err := auth.FindDomainID(cmd.Context(), identityClient, p.projectFlags.DomainNameOrID)
	if err == nil {
		domainName, err = auth.FindDomainName(cmd.Context(), identityClient, domainID)
	}
	if err != nil {
		return err
	}

	res := ratesProjects.List(cmd.Context(), limesRatesClient, domainID, ratesProjects.ReadOpts{
		Areas:    p.filterFlags.areas,
		Services: util.CastStringsTo[limes.ServiceType](p.filterFlags.services),
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get project reports")
	}

	if p.outputFmtFlags.format == core.OutputFormatJSON {
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

	projectFlags   projectFlags
	filterFlags    resourceFilterFlags
	outputFmtFlags resourceOutputFmtFlags
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
	projectShow.filterFlags.AddToCmd(cmd)
	projectShow.outputFmtFlags.AddToCmd(cmd)

	projectShow.Command = cmd
	return projectShow
}

func (p *projectShowCmd) Run(cmd *cobra.Command, args []string) error {
	nameOrID := ""
	if len(args) > 0 {
		nameOrID = args[0]
	}
	err := p.projectFlags.validateWithNameID(nameOrID)
	if err != nil {
		return err
	}

	outputOpts, err := p.outputFmtFlags.validate()
	if err != nil {
		return err
	}

	pInfo, err := auth.FindProject(cmd.Context(), identityClient, p.projectFlags.DomainNameOrID, nameOrID)
	if err != nil {
		return err
	}

	res := projects.Get(cmd.Context(), limesResourcesClient, pInfo.DomainID, pInfo.ID, projects.GetOpts{
		Areas:     p.filterFlags.areas,
		Services:  util.CastStringsTo[limes.ServiceType](p.filterFlags.services),
		Resources: util.CastStringsTo[limesresources.ResourceName](p.filterFlags.resources),
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get project report")
	}

	if p.outputFmtFlags.format == core.OutputFormatJSON {
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

	projectFlags   projectFlags
	filterFlags    rateFilterFlags
	outputFmtFlags rateOutputFmtFlags
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
	projectShowRates.filterFlags.AddToCmd(cmd)
	projectShowRates.outputFmtFlags.AddToCmd(cmd)

	projectShowRates.Command = cmd
	return projectShowRates
}

func (p *projectShowRatesCmd) Run(cmd *cobra.Command, args []string) error {
	nameOrID := ""
	if len(args) > 0 {
		nameOrID = args[0]
	}
	err := p.projectFlags.validateWithNameID(nameOrID)
	if err != nil {
		return err
	}

	outputOpts, err := p.outputFmtFlags.validate()
	if err != nil {
		return err
	}

	pInfo, err := auth.FindProject(cmd.Context(), identityClient, p.projectFlags.DomainNameOrID, nameOrID)
	if err != nil {
		return err
	}

	res := ratesProjects.Get(cmd.Context(), limesRatesClient, pInfo.DomainID, pInfo.ID, ratesProjects.ReadOpts{
		Areas:    p.filterFlags.areas,
		Services: util.CastStringsTo[limes.ServiceType](p.filterFlags.services),
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get project report")
	}

	if p.outputFmtFlags.format == core.OutputFormatJSON {
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
// Project sync.

type projectSyncCmd struct {
	*cobra.Command

	flags projectFlags
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
	projectSync.flags.AddToCmd(cmd)

	projectSync.Command = cmd
	return projectSync
}

func (p *projectSyncCmd) Run(cmd *cobra.Command, args []string) error {
	nameOrID := ""
	if len(args) > 0 {
		nameOrID = args[0]
	}
	err := p.flags.validateWithNameID(nameOrID)
	if err != nil {
		return err
	}

	pInfo, err := auth.FindProject(cmd.Context(), identityClient, p.flags.DomainNameOrID, nameOrID)
	if err != nil {
		return err
	}

	err = projects.Sync(cmd.Context(), limesResourcesClient, pInfo.DomainID, pInfo.ID).ExtractErr()
	if err != nil {
		return util.WrapError(err, "could not sync project")
	}

	return nil
}
