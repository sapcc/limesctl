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
	"github.com/sapcc/gophercloud-sapcc/resources/v1/clusters"
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
	return cmd
}

///////////////////////////////////////////////////////////////////////////////
// Cluster show.

type clusterShowCmd struct {
	*cobra.Command

	resourceFilterFlags
	resourceOutputFmtFlags
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
	clusterShow.resourceFilterFlags.AddToCmd(cmd)
	clusterShow.resourceOutputFmtFlags.AddToCmd(cmd)

	clusterShow.Command = cmd
	return clusterShow
}

func (c *clusterShowCmd) Run(_ *cobra.Command, _ []string) error {
	outputOpts, err := c.resourceOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	res := clusters.Get(limesResourcesClient, clusters.GetOpts{
		Areas:     c.areas,
		Services:  c.services,
		Resources: c.resources,
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get cluster report")
	}

	if c.format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesRep, err := res.Extract()
	if err != nil {
		return util.WrapError(err, "could not extract cluster report")
	}

	return writeReports(outputOpts, core.ClusterReport{ClusterReport: limesRep})
}
