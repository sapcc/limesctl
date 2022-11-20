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
	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/go-api-declarations/limes"
	"github.com/sapcc/gophercloud-sapcc/resources/v1/domains"
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
	cmd.AddCommand(newDomainSetCmd().Command)
	return cmd
}

///////////////////////////////////////////////////////////////////////////////
// Domain list.

type domainListCmd struct {
	*cobra.Command

	resourceFilterFlags
	resourceOutputFmtFlags
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
	domainList.resourceFilterFlags.AddToCmd(cmd)
	domainList.resourceOutputFmtFlags.AddToCmd(cmd)

	domainList.Command = cmd
	return domainList
}

func (d *domainListCmd) Run(_ *cobra.Command, _ []string) error {
	outputOpts, err := d.resourceOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	res := domains.List(limesResourcesClient, domains.ListOpts{
		Areas:     d.areas,
		Services:  d.services,
		Resources: d.resources,
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get domain reports")
	}

	if d.format == core.OutputFormatJSON {
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

	resourceFilterFlags
	resourceOutputFmtFlags
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
	domainShow.resourceFilterFlags.AddToCmd(cmd)
	domainShow.resourceOutputFmtFlags.AddToCmd(cmd)

	domainShow.Command = cmd
	return domainShow
}

func (d *domainShowCmd) Run(_ *cobra.Command, args []string) error {
	outputOpts, err := d.resourceOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	nameOrID := ""
	if len(args) > 0 {
		nameOrID = args[0]
	}
	domainID, err := auth.FindDomainID(identityClient, nameOrID)
	if err != nil {
		return err
	}

	res := domains.Get(limesResourcesClient, domainID, domains.GetOpts{
		Areas:     d.areas,
		Services:  d.services,
		Resources: d.resources,
	})
	if res.Err != nil {
		return util.WrapError(res.Err, "could not get domain report")
	}

	if d.format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesRep, err := res.Extract()
	if err != nil {
		return util.WrapError(err, "could not extract domain report")
	}

	return writeReports(outputOpts, core.DomainReport{DomainReport: limesRep})
}

///////////////////////////////////////////////////////////////////////////////
// Domain set.

type domainSetCmd struct {
	*cobra.Command

	quotas []string
}

func newDomainSetCmd() *domainSetCmd {
	domainSet := &domainSetCmd{}
	cmd := &cobra.Command{
		Use:   "set [name or ID]",
		Short: "Change resource quota values for a specific domain",
		Long: `Change resource quota values for a specific domain.

For relative quota adjustment, use one of the following operators: [+=, -=, *=, /=].

This command requires a cloud-admin token.`,
		Example: makeExamplesString([]string{
			`limesctl domain set --quotas="compute/cores=200,compute/ram=50GiB"`,
			`limesctl domain set -q compute/cores=200 -q compute/ram=50GiB (you can give the flag multiple times)`,
			`limesctl domain set -q compute/cores=200 -q compute/ram=50GiB (you can give the flag multiple times)`,
			`limesctl domain set -q compute/cores*=2 -q compute/ram+=10GiB`,
			`limesctl domain set -q object-store/capacity=1TiB (you can also use a unit other than the service's default, e.g. object-store uses 'B' by default but we use 'TiB' here)`,
			`limesctl domain set -q object-store/capacity-=0.25TiB (fractional values are also possible)`,
		}),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: authWithLimesResources,
		RunE:    domainSet.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	cmd.Flags().StringSliceVarP(&domainSet.quotas, "quotas", "q", nil, "new quota values (comma separated list)")
	cmd.MarkFlagRequired("quotas")

	domainSet.Command = cmd
	return domainSet
}

func (d *domainSetCmd) Run(_ *cobra.Command, args []string) error {
	nameOrID := ""
	if len(args) > 0 {
		nameOrID = args[0]
	}
	domainID, err := auth.FindDomainID(identityClient, nameOrID)
	if err != nil {
		return err
	}

	resQuotas, err := getDomainResourceQuotas(limesResourcesClient, domainID)
	if err != nil {
		return util.WrapError(err, "could not get default units")
	}
	qc, err := parseToQuotaRequest(resQuotas, d.quotas)
	if err != nil {
		return util.WrapError(err, "could not parse quota values")
	}

	err = domains.Update(limesResourcesClient, domainID, domains.UpdateOpts{
		Services: qc,
	}).ExtractErr()
	if err != nil {
		return util.WrapError(err, "could not set new quotas for domain")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Helper functions.

func getDomainResourceQuotas(limesClient *gophercloud.ServiceClient, domainID string) (resourceQuotas, error) {
	rep, err := domains.Get(limesClient, domainID, domains.GetOpts{}).Extract()
	if err != nil {
		return nil, err
	}

	result := make(resourceQuotas)
	for srv, srvReport := range rep.Services {
		for res, resReport := range srvReport.Resources {
			if _, ok := result[srv]; !ok {
				result[srv] = make(map[string]limes.ValueWithUnit)
			}
			var val uint64
			if resReport.DomainQuota != nil {
				val = *resReport.DomainQuota
			}
			result[srv][res] = limes.ValueWithUnit{Value: val, Unit: resReport.ResourceInfo.Unit}
		}
	}
	return result, nil
}
