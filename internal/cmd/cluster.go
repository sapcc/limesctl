// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/sapcc/go-api-declarations/limes"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
	ratesClusters "github.com/sapcc/gophercloud-sapcc/v2/rates/v1/clusters"
	"github.com/sapcc/gophercloud-sapcc/v2/resources/v1/clusters"
	"github.com/spf13/cobra"

	"github.com/sapcc/limesctl/v3/internal/core"
	"github.com/sapcc/limesctl/v3/internal/util"
)

func newClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Do some action at cluster level",
		Args:  cobra.NoArgs,
	}
	doNotSortFlags(cmd)
	cmd.AddCommand(newClusterShowCmd().Command)
	cmd.AddCommand(newClusterShowRatesCmd().Command)
	cmd.AddCommand(newClusterMailTemplateRenderer().Command)
	return cmd
}

///////////////////////////////////////////////////////////////////////////////
// Cluster show.

type clusterShowCmd struct {
	*cobra.Command

	filterFlags    resourceFilterFlags
	outputFmtFlags resourceOutputFmtFlags
}

func newClusterShowCmd() *clusterShowCmd {
	clusterShow := &clusterShowCmd{}
	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Display resource usage data for cluster",
		Long:    "Display resource usage data for cluster. This command requires a cloud-admin token.",
		Args:    cobra.NoArgs,
		PreRunE: authWithLimesResources,
		RunE:    clusterShow.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	clusterShow.filterFlags.AddToCmd(cmd)
	clusterShow.outputFmtFlags.AddToCmd(cmd)

	clusterShow.Command = cmd
	return clusterShow
}

// Run is called by Cobra when this command is executed.
func (c *clusterShowCmd) Run(cmd *cobra.Command, _ []string) error {
	outputOpts, err := c.outputFmtFlags.validate()
	if err != nil {
		return err
	}

	res := clusters.Get(cmd.Context(), limesResourcesClient, clusters.GetOpts{
		Areas:     c.filterFlags.areas,
		Services:  util.CastStringsTo[limes.ServiceType](c.filterFlags.services),
		Resources: util.CastStringsTo[limesresources.ResourceName](c.filterFlags.resources),
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get cluster report")
	}

	if c.outputFmtFlags.format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesRep, err := res.Extract()
	if err != nil {
		return util.WrapError(err, "could not extract cluster report")
	}

	return writeReports(outputOpts, core.ClusterReport{ClusterReport: limesRep})
}

///////////////////////////////////////////////////////////////////////////////
// Cluster show rates.

type clusterShowRatesCmd struct {
	*cobra.Command

	filterFlags    rateFilterFlags
	outputFmtFlags rateOutputFmtFlags
}

func newClusterShowRatesCmd() *clusterShowRatesCmd {
	clusterShowRates := &clusterShowRatesCmd{}
	cmd := &cobra.Command{
		Use:     "show-rates",
		Short:   "Display global rate limits for the cluster",
		Long:    "Display global rate limits for the cluster level. These rate limits apply to all users in aggregate.",
		Args:    cobra.NoArgs,
		PreRunE: authWithLimesRates,
		RunE:    clusterShowRates.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	clusterShowRates.filterFlags.AddToCmd(cmd)
	clusterShowRates.outputFmtFlags.AddToCmd(cmd)

	clusterShowRates.Command = cmd
	return clusterShowRates
}

// Run is called by Cobra when this command is executed.
func (c *clusterShowRatesCmd) Run(cmd *cobra.Command, args []string) error {
	outputOpts, err := c.outputFmtFlags.validate()
	if err != nil {
		return err
	}

	res := ratesClusters.Get(cmd.Context(), limesRatesClient, ratesClusters.GetOpts{
		Areas:    c.filterFlags.areas,
		Services: util.CastStringsTo[limes.ServiceType](c.filterFlags.services),
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get cluster report")
	}

	if c.outputFmtFlags.format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesRep, err := res.Extract()
	if err != nil {
		return util.WrapError(err, "could not extract cluster report")
	}

	return writeReports(outputOpts, core.ClusterRatesReport{ClusterReport: limesRep})
}

type clusterMailTemplate struct {
	*cobra.Command

	jsonOutput             bool
	confirmedCommitments   bool
	expiringCommitments    bool
	transferredCommitments bool
}

type mailTemplates struct {
	ConfirmedCommitments   string `json:"confirmed_commitments"`
	ExpiringCommitments    string `json:"expiring_commitments"`
	TransferredCommitments string `json:"transferred_commitments"`
}

func newClusterMailTemplateRenderer() *clusterMailTemplate {
	clusterMail := &clusterMailTemplate{}
	cmd := &cobra.Command{
		Use:     "show-mail-templates",
		Short:   "Display configured mail templates.",
		Long:    "Display configured mail templates. Can be used to check the validity of the configured templates.",
		Args:    cobra.NoArgs,
		PreRunE: authWithLimesAdmin,
		RunE:    clusterMail.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	cmd.Flags().BoolVar(&clusterMail.jsonOutput, "json", false, "output as JSON (with escaped HTML characters)")
	cmd.Flags().BoolVar(&clusterMail.confirmedCommitments, "confirmed", false, "show only confirmed commitments template")
	cmd.Flags().BoolVar(&clusterMail.expiringCommitments, "expiring", false, "show only expiring commitments template")
	cmd.Flags().BoolVar(&clusterMail.transferredCommitments, "transferred", false, "show only transferred commitments template")
	clusterMail.Command = cmd
	return clusterMail
}

// Run is called by Cobra when this command is executed.
func (c *clusterMailTemplate) Run(cmd *cobra.Command, args []string) error {
	var res gophercloud.Result
	resp, err := limesAdminClient.Get(cmd.Context(), limesAdminClient.ServiceURL("mail/render"), &res.Body, nil) //nolint:bodyclose
	_, _, res.Err = gophercloud.ParseResponse(resp, err)
	if res.Err != nil {
		return util.WrapError(res.Err, "could not fetch mail template request body from limes")
	}

	templates := &mailTemplates{}
	err = res.ExtractInto(templates)
	if err != nil {
		return util.WrapError(err, "could not extract mail templates")
	}

	switch {
	case c.confirmedCommitments:
		if c.jsonOutput {
			return writeJSON(templates.ConfirmedCommitments)
		}
		fmt.Println(templates.ConfirmedCommitments)
		return nil
	case c.expiringCommitments:
		if c.jsonOutput {
			return writeJSON(templates.ExpiringCommitments)
		}
		fmt.Println(templates.ExpiringCommitments)
		return nil
	case c.transferredCommitments:
		if c.jsonOutput {
			return writeJSON(templates.TransferredCommitments)
		}
		fmt.Println(templates.TransferredCommitments)
		return nil
	default:
		if c.jsonOutput {
			return writeJSON(res.Body)
		}
		fmt.Println(templates.ConfirmedCommitments)
		fmt.Println(templates.ExpiringCommitments)
		fmt.Println(templates.TransferredCommitments)
		return nil
	}
}
