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

type ProjectCmd struct {
	Cluster string `short:"c" help:"Cluster ID. When this option is used, the domain/project must be identified by ID. Specifying a domain/project name will not work."`

	// TODO: check if domain name will work?
	Domain string `short:"d" help:"Domain ID. When this option is used, the domain/project must be identified by ID. Specifying a domain/project name will not work."`

	List projectListCmd `cmd:"" help:"Display data for all the projects. Requires a domain-admin token."`
	Show projectShowCmd `cmd:"" help:"Display data for a specific project. Requires project member permissions."`
	Set  projectSetCmd  `cmd:"" help:"Change quota values for a specific project. Requires a domain-admin token."`
	Sync projectSyncCmd `cmd:"" help:"Schedule a sync job that pulls quota and usage data for this project from the backing services into Limes' local database. Requires project-admin token."`
}

type projectListCmd struct {
	RequestFilterFlags
	OutputFormatFlags
}

type projectShowCmd struct {
	RequestFilterFlags
	OutputFormatFlags

	NameOrID string `arg:"" optional:"" help:"Name or ID of the project (leave empty for current project)."`
}

type projectSetCmd struct {
	Quotas   []string `help:"New quotas values. Example: service/resource=120GiB."`
	NameOrID string   `arg:"" optional:"" help:"Name or ID of the project (leave empty for current project)."`
}

type projectSyncCmd struct {
	NameOrID string `arg:"" required:"" help:"Name or ID of the project (leave empty for current project)."`
}
