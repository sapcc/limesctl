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
	"strings"

	"github.com/alecthomas/kong"

	"github.com/sapcc/limesctl/internal/cmd"
	"github.com/sapcc/limesctl/internal/core"
)

// This info identifies a specific build of the app. It is set at compile time.
var (
	version = "dev"
	commit  = "unknown"
	date    = "now"
)

type rootCmd struct {
	cmd.Globals

	Cluster cmd.ClusterCmd `cmd:"" help:"Do some action on cluster(s)."`
	Domain  cmd.DomainCmd  `cmd:"" help:"Do some action on domain(s)."`
	Project cmd.ProjectCmd `cmd:"" help:"Do some action on project(s)."`
}

func main() {
	var cli rootCmd
	ctx := kong.Parse(&cli,
		kong.Name("limesctl"),
		kong.Description("Command-line client for Limes."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Bind(cmd.VersionFlag{
			Version:       version,
			GitCommitHash: commit,
			BuildDate:     date,
		}),
		kong.Vars{"outputFormats": outputFormats()},
	)

	clients, err := cli.Globals.Authenticate()
	if err == nil {
		err = ctx.Run(clients)
	}
	if err != nil {
		ctx.FatalIfErrorf(err)
	}
}

func outputFormats() string {
	f := []string{
		string(core.OutputFormatTable),
		string(core.OutputFormatCSV),
		string(core.OutputFormatJSON),
	}
	return strings.Join(f, ",")
}
