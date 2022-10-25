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
	"github.com/pkg/errors"
	"github.com/sapcc/gophercloud-sapcc/resources/v1/clusters"

	"github.com/sapcc/limesctl/v3/internal/core"
)

// clusterCmd contains the command-line structure for the cluster command.
type clusterCmd struct {
	Show clusterShowCmd `cmd:"" help:"Display resource usage data for the cluster. Requires a cloud-admin token."`
}

///////////////////////////////////////////////////////////////////////////////
// Cluster show.

type clusterShowCmd struct {
	resourceFilterFlags
	resourceOutputFmtFlags

	ID string `arg:"" optional:"" help:"ID of the cluster (leave empty for current cluster)."`
}

func (c *clusterShowCmd) Run(clients *ServiceClients) error {
	outputOpts, err := c.resourceOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	if c.ID == "" {
		c.ID = "current"
	}
	res := clusters.Get(clients.limesResources, clusters.GetOpts{
		Areas:     c.Areas,
		Services:  c.Services,
		Resources: c.Resources,
	})
	if res.Err != nil {
		return errors.Wrap(res.Err, "could not get cluster report")
	}

	if c.Format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesRep, err := res.Extract()
	if err != nil {
		return errors.Wrap(err, "could not extract cluster report")
	}

	return writeReports(outputOpts, core.ClusterReport{ClusterReport: limesRep})
}
