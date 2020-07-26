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

type DomainCmd struct {
	Cluster string `short:"c" help:"Cluster ID. When this option is used, the domain must be identified by ID. Specifying a domain name will not work."`

	List domainListCmd `cmd:"" help:"Display data for all the domains. Requires a cloud-admin token."`
	Show domainShowCmd `cmd:"" help:"Display data for a specific domain. Requires a domain-admin token."`
	Set  domainSetCmd  `cmd:"" help:"Change quota values for a specific domain. Requires a cloud-admin token."`
}

type domainListCmd struct {
	RequestFilterFlags
	OutputFormatFlags
}

type domainShowCmd struct {
	RequestFilterFlags
	OutputFormatFlags

	NameOrID string `arg:"" optional:"" help:"Name or ID of the domain (leave empty for current domain)."`
}

type domainSetCmd struct {
	Quotas   []string `help:"New quotas values. Example: service/resource=120GiB."`
	NameOrID string   `arg:"" optional:"" help:"Name or ID of the domain (leave empty for current domain)."`
}
