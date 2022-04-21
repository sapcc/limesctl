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
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/pkg/errors"
	"github.com/sapcc/go-api-declarations/limes"
	"github.com/sapcc/gophercloud-sapcc/resources/v1/projects"

	"github.com/sapcc/limesctl/internal/auth"
	"github.com/sapcc/limesctl/internal/core"
)

//nolint:lll
// projectCmd contains the command-line structure for the project command.
type projectCmd struct {
	List      projectListCmd      `cmd:"" help:"Display resource usage data for all the projects in a domain. Requires a domain-admin token."`
	ListRates projectListRatesCmd `cmd:"" help:"Display rate limits for all the projects in a domain. Requires a domain-admin token."`
	Show      projectShowCmd      `cmd:"" help:"Display resource usage data for a specific project. Requires project member permissions."`
	ShowRates projectShowRatesCmd `cmd:"" help:"Display rate limits for a specific project. Requires project member permissions."`
	Set       projectSetCmd       `cmd:"" help:"Change resource quota values for a specific project. Requires a domain-admin token."`
	Sync      projectSyncCmd      `cmd:"" help:"Schedule a sync job that pulls quota and usage data for a specific project from the backing services into Limes' local database. Requires a project-admin token."`
}

type projectFlags struct {
	DomainNameOrID string `short:"d" name:"domain" help:"Name or ID of the domain."`
}

func (pf *projectFlags) validateWithNameID(nameOrID string) error {
	if pf.DomainNameOrID != "" && nameOrID == "" {
		return errors.New("Project name or ID is required when using the '--domain' flag")
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Project list.

type projectListCmd struct {
	projectFlags
	resourceFilterFlags
	resourceOutputFmtFlags
}

func (p *projectListCmd) Run(clients *ServiceClients) error {
	outputOpts, err := p.resourceOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	domainName := ""
	domainID, err := auth.FindDomainID(clients.identity, p.DomainNameOrID)
	if err == nil {
		domainName, err = auth.FindDomainName(clients.identity, domainID)
	}
	if err != nil {
		return err
	}

	res := projects.List(clients.limes, domainID, projects.ListOpts{
		Areas:     p.Areas,
		Services:  p.Services,
		Resources: p.Resources,
	})
	if res.Err != nil {
		return errors.Wrap(res.Err, "could not get project reports")
	}

	if p.Format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesReps, err := res.ExtractProjects()
	if err != nil {
		return errors.Wrap(err, "could not extract project reports")
	}

	return writeReports(outputOpts,
		core.LimesProjectsToReportRenderer(limesReps, domainID, domainName, false)...)
}

///////////////////////////////////////////////////////////////////////////////
// Project list rates.

type projectListRatesCmd struct {
	projectFlags
	rateFilterFlags
	rateOutputFmtFlags
}

func (p *projectListRatesCmd) Run(clients *ServiceClients) error {
	outputOpts, err := p.rateOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	domainName := ""
	domainID, err := auth.FindDomainID(clients.identity, p.DomainNameOrID)
	if err == nil {
		domainName, err = auth.FindDomainName(clients.identity, domainID)
	}
	if err != nil {
		return err
	}

	res := projects.List(clients.limes, domainID, projects.ListOpts{
		Services: p.Services,
		Rates:    projects.OnlyRates,
	})
	if res.Err != nil {
		return errors.Wrap(res.Err, "could not get project reports")
	}

	if p.Format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesReps, err := res.ExtractProjects()
	if err != nil {
		return errors.Wrap(err, "could not extract project reports")
	}

	return writeReports(outputOpts,
		core.LimesProjectsToReportRenderer(limesReps, domainID, domainName, true)...)
}

///////////////////////////////////////////////////////////////////////////////
// Project show.

type projectShowCmd struct {
	projectFlags
	resourceFilterFlags
	resourceOutputFmtFlags

	NameOrID string `arg:"" optional:"" help:"Name or ID of the project. Required if using '--domain' flag."`
}

// Validate implements the kong.Validatable interface.
func (p *projectShowCmd) Validate() error {
	return p.projectFlags.validateWithNameID(p.NameOrID)
}

func (p *projectShowCmd) Run(clients *ServiceClients) error {
	outputOpts, err := p.resourceOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	pInfo, err := auth.FindProject(clients.identity, p.DomainNameOrID, p.NameOrID)
	if err != nil {
		return err
	}

	res := projects.Get(clients.limes, pInfo.DomainID, pInfo.ID, projects.GetOpts{
		Areas:     p.Areas,
		Services:  p.Services,
		Resources: p.Resources,
	})
	if res.Err != nil {
		return errors.Wrap(res.Err, "could not get project report")
	}

	if p.Format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesRep, err := res.Extract()
	if err != nil {
		return errors.Wrap(err, "could not extract project report")
	}

	return writeReports(outputOpts, core.ProjectReport{
		ProjectReport: limesRep,
		DomainID:      pInfo.DomainID,
		DomainName:    pInfo.DomainName,
	})
}

///////////////////////////////////////////////////////////////////////////////
// Project show rates.

type projectShowRatesCmd struct {
	projectFlags
	rateFilterFlags
	rateOutputFmtFlags

	NameOrID string `arg:"" optional:"" help:"Name or ID of the project. Required if using '--domain' flag."`
}

// Validate implements the kong.Validatable interface.
func (p *projectShowRatesCmd) Validate() error {
	return p.projectFlags.validateWithNameID(p.NameOrID)
}

func (p *projectShowRatesCmd) Run(clients *ServiceClients) error {
	outputOpts, err := p.rateOutputFmtFlags.validate()
	if err != nil {
		return err
	}

	pInfo, err := auth.FindProject(clients.identity, p.DomainNameOrID, p.NameOrID)
	if err != nil {
		return err
	}

	res := projects.Get(clients.limes, pInfo.DomainID, pInfo.ID, projects.GetOpts{
		Services: p.Services,
		Rates:    projects.OnlyRates,
	})
	if res.Err != nil {
		return errors.Wrap(res.Err, "could not get project report")
	}

	if p.Format == core.OutputFormatJSON {
		return writeJSON(res.Body)
	}

	limesRep, err := res.Extract()
	if err != nil {
		return errors.Wrap(err, "could not extract project report")
	}

	return writeReports(outputOpts, core.ProjectReport{
		ProjectReport: limesRep,
		HasRatesOnly:  true,
		DomainID:      pInfo.DomainID,
		DomainName:    pInfo.DomainName,
	})
}

///////////////////////////////////////////////////////////////////////////////
// Project set.

//nolint:lll
type projectSetCmd struct {
	projectFlags
	Quotas []string `short:"q" sep:"," help:"New quotas values. For relative quota adjustment, use one of the following operators: [+=, -=, *=, /=]. Example: service/resource=120GiB."`

	NameOrID string `arg:"" optional:"" help:"Name or ID of the project. Required if using '--domain' flag."`
}

// Validate implements the kong.Validatable interface.
func (p *projectSetCmd) Validate() error {
	return p.projectFlags.validateWithNameID(p.NameOrID)
}

func (p *projectSetCmd) Run(clients *ServiceClients) error {
	pInfo, err := auth.FindProject(clients.identity, p.DomainNameOrID, p.NameOrID)
	if err != nil {
		return err
	}

	resQuotas, err := getProjectResourceQuotas(clients.limes, pInfo)
	if err != nil {
		return errors.Wrap(err, "could not get default units")
	}
	qc, err := parseToQuotaRequest(resQuotas, p.Quotas)
	if err != nil {
		return errors.Wrap(err, "could not parse quota values")
	}

	warn, err := projects.Update(clients.limes, pInfo.DomainID, pInfo.ID, projects.UpdateOpts{
		Services: qc,
	}).Extract()
	if err != nil {
		return errors.Wrap(err, "could not set new quotas for project")
	}

	if warn != nil {
		fmt.Println(string(warn))
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Project sync.

type projectSyncCmd struct {
	projectFlags

	NameOrID string `arg:"" required:"" help:"Name or ID of the project."`
}

// Validate implements the kong.Validatable interface.
func (p *projectSyncCmd) Validate() error {
	return p.projectFlags.validateWithNameID(p.NameOrID)
}

func (p *projectSyncCmd) Run(clients *ServiceClients) error {
	pInfo, err := auth.FindProject(clients.identity, p.DomainNameOrID, p.NameOrID)
	if err != nil {
		return err
	}

	err = projects.Sync(clients.limes, pInfo.DomainID, pInfo.ID).ExtractErr()
	if err != nil {
		return errors.Wrap(err, "could not sync project")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Helper functions.

func getProjectResourceQuotas(limesClient *gophercloud.ServiceClient, pInfo *auth.ProjectInfo) (resourceQuotas, error) {
	rep, err := projects.Get(limesClient, pInfo.DomainID, pInfo.ID, projects.GetOpts{}).Extract()
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
			if resReport.Quota != nil {
				val = *resReport.Quota
			}
			result[srv][res] = limes.ValueWithUnit{Value: val, Unit: resReport.ResourceInfo.Unit}
		}
	}
	return result, nil
}
