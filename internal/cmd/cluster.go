// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
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
