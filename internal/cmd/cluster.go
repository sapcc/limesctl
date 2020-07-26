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

type ClusterCmd struct {
	List clusterListCmd `cmd:"" help:"Display data for all the clusters. Requires a cloud-admin token."`
	Show clusterShowCmd `cmd:"" help:"Display data for a specific cluster. Requires a cloud-admin token."`
	Set  clusterSetCmd  `cmd:"" help:"Change capacity values for a specific cluster. Requires a cloud-admin token."`
}

type clusterListCmd struct {
	RequestFilterFlags
	OutputFormatFlags
}

type clusterShowCmd struct {
	RequestFilterFlags
	OutputFormatFlags

	NameOrID string `arg:"" optional:"" help:"Name or ID of the cluster (leave empty for current cluster)."`
}

type clusterSetCmd struct {
	Capacities []string `help:"New capacity values. Example: service/resource=256TiB:'Tis a comment' (comment is optional)."`
	NameOrID   string   `arg:"" optional:"" help:"Name or ID of the cluster (leave empty for current cluster)."`
}
