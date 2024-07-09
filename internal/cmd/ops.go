/*******************************************************************************
*
* Copyright 2024 SAP SE
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

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/sapcc/go-api-declarations/limes"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
	"github.com/sapcc/gophercloud-sapcc/v2/resources/v1/projects"
	"github.com/spf13/cobra"

	"github.com/sapcc/limesctl/v3/internal/auth"
	"github.com/sapcc/limesctl/v3/internal/util"
)

func newOpsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ops",
		Short: "Toolbox for Limes operators (end users do not need this)",
		Args:  cobra.NoArgs,
	}
	// Flags
	doNotSortFlags(cmd)
	// Subcommands
	cmd.AddCommand(newOpsValidateQuotaOverridesCmd())
	return cmd
}

func newOpsValidateQuotaOverridesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate-quota-overrides path",
		Short: "Validate a quota-overrides.json file for usage with an existing Limes instance",
		Long: `Validate a quota-overrides.json file for usage with an existing Limes instance.

Requies project member permissions.`,
		Args:    cobra.ExactArgs(1),
		PreRunE: authWithLimesResources,
		RunE:    runValidateQuotaOverrides,
	}
}

func runValidateQuotaOverrides(cmd *cobra.Command, args []string) error {
	buf, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	// get resource report for the project from the current token scope
	pInfo, err := auth.FindProject(cmd.Context(), identityClient, "", "")
	if err != nil {
		return err
	}
	report, err := projects.Get(cmd.Context(), limesResourcesClient, pInfo.DomainID, pInfo.ID, projects.GetOpts{}).Extract()
	if err != nil {
		return util.WrapError(err, "could not get project report")
	}

	// parse quota-overrides.json file using the existing report as
	// a hint for which services/resources exist with which units
	getUnit := func(serviceType limes.ServiceType, resourceName limesresources.ResourceName) (limes.Unit, error) {
		srvReport := report.Services[serviceType]
		if srvReport == nil {
			return limes.UnitUnspecified, fmt.Errorf("%q is not a valid service", serviceType)
		}
		resReport := srvReport.Resources[resourceName]
		if resReport == nil {
			fullResourceName := fmt.Sprintf("%s/%s", serviceType, resourceName)
			return limes.UnitUnspecified, fmt.Errorf("%q is not a valid resource", fullResourceName)
		}
		if resReport.Quota == nil {
			return limes.UnitUnspecified, fmt.Errorf("%s/%s does not track quota", serviceType, resourceName)
		}
		return resReport.Unit, nil
	}
	_, errs := limesresources.ParseQuotaOverrides(buf, getUnit)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		}
		return errors.New("validation failed")
	}

	return nil
}
