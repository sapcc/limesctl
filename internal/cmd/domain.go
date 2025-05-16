// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/sapcc/go-api-declarations/limes"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
	"github.com/sapcc/gophercloud-sapcc/v2/resources/v1/domains"
	"github.com/spf13/cobra"

	"github.com/sapcc/limesctl/v3/internal/auth"
	"github.com/sapcc/limesctl/v3/internal/core"
	"github.com/sapcc/limesctl/v3/internal/util"
)

func newDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Do some action at domain level",
		Args:  cobra.NoArgs,
	}
	// Flags
	doNotSortFlags(cmd)
	// Subcommands
	cmd.AddCommand(newDomainListCmd().Command)
	cmd.AddCommand(newDomainShowCmd().Command)
	return cmd
}

///////////////////////////////////////////////////////////////////////////////
// Domain list.

type domainListCmd struct {
	*cobra.Command

	filterFlags    resourceFilterFlags
	outputFmtFlags resourceOutputFmtFlags
}

func newDomainListCmd() *domainListCmd {
	domainList := &domainListCmd{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Display resource usage data for all domains",
		Long: `Display resource usage data for all domains. This command requires a
cloud-admin token.`,
		Args:    cobra.NoArgs,
		PreRunE: authWithLimesResources,
		RunE:    domainList.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	domainList.filterFlags.AddToCmd(cmd)
	domainList.outputFmtFlags.AddToCmd(cmd)

	domainList.Command = cmd
	return domainList
}

func (d *domainListCmd) Run(cmd *cobra.Command, _ []string) error {
	outputOpts, err := d.outputFmtFlags.validate()
	if err != nil {
		return err
	}

	res := domains.List(cmd.Context(), limesResourcesClient, domains.ListOpts{
		Areas:     d.filterFlags.areas,
		Services:  util.CastStringsTo[limes.ServiceType](d.filterFlags.services),
		Resources: util.CastStringsTo[limesresources.ResourceName](d.filterFlags.resources),
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get domain reports")
	}

	if d.outputFmtFlags.format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesReps, err := res.ExtractDomains()
	if err != nil {
		return util.WrapError(err, "could not extract domain reports")
	}

	return writeReports(outputOpts, core.LimesDomainsToReportRenderer(limesReps)...)
}

///////////////////////////////////////////////////////////////////////////////
// Domain show.

type domainShowCmd struct {
	*cobra.Command

	filterFlags    resourceFilterFlags
	outputFmtFlags resourceOutputFmtFlags
}

func newDomainShowCmd() *domainShowCmd {
	domainShow := &domainShowCmd{}
	cmd := &cobra.Command{
		Use:   "show [name or ID]",
		Short: "Display resource usage data for a specific domain",
		Long: `Display resource usage data for a specific domain. This command requires a
domain-admin token.`,
		Args:    cobra.MaximumNArgs(1),
		PreRunE: authWithLimesResources,
		RunE:    domainShow.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	domainShow.filterFlags.AddToCmd(cmd)
	domainShow.outputFmtFlags.AddToCmd(cmd)

	domainShow.Command = cmd
	return domainShow
}

func (d *domainShowCmd) Run(cmd *cobra.Command, args []string) error {
	outputOpts, err := d.outputFmtFlags.validate()
	if err != nil {
		return err
	}

	nameOrID := ""
	if len(args) > 0 {
		nameOrID = args[0]
	}
	domainID, err := auth.FindDomainID(cmd.Context(), identityClient, nameOrID)
	if err != nil {
		return err
	}

	res := domains.Get(cmd.Context(), limesResourcesClient, domainID, domains.GetOpts{
		Areas:     d.filterFlags.areas,
		Services:  util.CastStringsTo[limes.ServiceType](d.filterFlags.services),
		Resources: util.CastStringsTo[limesresources.ResourceName](d.filterFlags.resources),
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get domain report")
	}

	if d.outputFmtFlags.format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesRep, err := res.Extract()
	if err != nil {
		return util.WrapError(err, "could not extract domain report")
	}

	return writeReports(outputOpts, core.DomainReport{DomainReport: limesRep})
}
