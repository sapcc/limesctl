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
	"github.com/pkg/errors"
	"github.com/sapcc/gophercloud-sapcc/resources/v1/domains"
	"github.com/sapcc/limes"

	"github.com/sapcc/limesctl/internal/auth"
	"github.com/sapcc/limesctl/internal/core"
)

// domainCmd contains the command-line structure for the domain command.
type domainCmd struct {
	List domainListCmd `cmd:"" help:"Display resource usage data for all the domains. Requires a cloud-admin token."`
	Show domainShowCmd `cmd:"" help:"Display resource usage data for a specific domain. Requires a domain-admin token."`
	Set  domainSetCmd  `cmd:"" help:"Change resource quota values for a specific domain. Requires a cloud-admin token."`
}

//nolint:lll
type domainClusterFlag struct {
	ClusterID string `short:"c" name:"cluster" help:"Cluster ID. When this option is used, the domain must be identified by ID (name won't work)."`
}

///////////////////////////////////////////////////////////////////////////////
// Domain list.

type domainListCmd struct {
	domainClusterFlag
	resourceFilterFlags
	resourceOutputFmtFlags
}

func (d *domainListCmd) Run(clients *ServiceClients) error {
	outputOpts, err := d.resourceOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	res := domains.List(clients.limes, domains.ListOpts{
		Cluster:   d.ClusterID,
		Areas:     d.Areas,
		Services:  d.Services,
		Resources: d.Resources,
	})
	if res.Err != nil {
		return errors.Wrap(res.Err, "could not get domain reports")
	}

	if d.Format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesReps, err := res.ExtractDomains()
	if err != nil {
		return errors.Wrap(err, "could not extract domain reports")
	}

	return writeReports(outputOpts, core.LimesDomainsToReportRenderer(limesReps)...)
}

///////////////////////////////////////////////////////////////////////////////
// Domain show.

type domainShowCmd struct {
	domainClusterFlag
	resourceFilterFlags
	resourceOutputFmtFlags

	NameOrID string `arg:"" optional:"" help:"Name or ID of the domain. Required if using '--cluster' flag."`
}

// Validate implements the kong.Validatable interface.
func (d *domainShowCmd) Validate() error {
	if d.ClusterID != "" && d.NameOrID == "" {
		return errors.New("Domain ID is required when using the '--cluster' flag")
	}
	return nil
}

func (d *domainShowCmd) Run(clients *ServiceClients) error {
	outputOpts, err := d.resourceOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	domainID := d.NameOrID
	if d.ClusterID == "" {
		var err error
		domainID, err = auth.FindDomainID(clients.identity, d.NameOrID)
		if err != nil {
			return err
		}
	}

	res := domains.Get(clients.limes, domainID, domains.GetOpts{
		Cluster:   d.ClusterID,
		Areas:     d.Areas,
		Services:  d.Services,
		Resources: d.Resources,
	})
	if res.Err != nil {
		return errors.Wrap(res.Err, "could not get domain report")
	}

	if d.Format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesRep, err := res.Extract()
	if err != nil {
		return errors.Wrap(err, "could not extract domain report")
	}

	return writeReports(outputOpts, core.DomainReport{DomainReport: limesRep})
}

///////////////////////////////////////////////////////////////////////////////
// Domain set.

//nolint:lll
type domainSetCmd struct {
	domainClusterFlag
	Quotas []string `short:"q" sep:"," help:"New quotas values. For relative quota adjustment, use one of the following operators: [+=, -=, *=, /=]. Example: service/resource=10GiB."`

	NameOrID string `arg:"" optional:"" help:"Name or ID of the domain. Required if using '--cluster' flag."`
}

// Validate implements the kong.Validatable interface.
func (d *domainSetCmd) Validate() error {
	if d.ClusterID != "" && d.NameOrID == "" {
		return errors.New("Domain ID is required when using the '--cluster' flag")
	}
	return nil
}

func (d *domainSetCmd) Run(clients *ServiceClients) error {
	domainID := d.NameOrID
	if d.ClusterID == "" {
		var err error
		domainID, err = auth.FindDomainID(clients.identity, d.NameOrID)
		if err != nil {
			return err
		}
	}

	resQuotas, err := getDomainResourceQuotas(clients.limes, d.ClusterID, domainID)
	if err != nil {
		return errors.Wrap(err, "could not get default units")
	}
	qc, err := parseToQuotaRequest(resQuotas, d.Quotas)
	if err != nil {
		return errors.Wrap(err, "could not parse quota values")
	}

	err = domains.Update(clients.limes, domainID, domains.UpdateOpts{
		Cluster:  d.ClusterID,
		Services: qc,
	}).ExtractErr()
	if err != nil {
		return errors.Wrap(err, "could not set new quotas for domain")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Helper functions.

func getDomainResourceQuotas(limesClient *gophercloud.ServiceClient, clusterID, domainID string) (resourceQuotas, error) {
	rep, err := domains.Get(limesClient, domainID, domains.GetOpts{
		Cluster: clusterID,
	}).Extract()
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
