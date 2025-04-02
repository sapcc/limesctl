// Copyright 2024 SAP SE
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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/sapcc/go-api-declarations/liquid"
	"github.com/sapcc/go-bits/liquidapi"
	"github.com/spf13/cobra"

	"github.com/sapcc/limesctl/v3/internal/util"
)

func prettyPrint(obj any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(obj)
}

func newLiquidCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid",
		Short: "Execute requests against a LIQUID API",
		Args:  cobra.NoArgs,
	}
	doNotSortFlags(cmd)
	cmd.AddCommand(newLiquidServiceInfoCmd().Command)
	cmd.AddCommand(newLiquidReportCapacityCmd().Command)
	cmd.AddCommand(newLiquidReportUsageCmd().Command)
	cmd.AddCommand(newLiquidSetQuotaCmd().Command)
	return cmd
}

///////////////////////////////////////////////////////////////////////////////
// LIQUID Service Info

type liquidServiceInfoCmd struct {
	*cobra.Command

	flags liquidOperationFlags
}

func newLiquidServiceInfoCmd() *liquidServiceInfoCmd {
	liquidServiceInfo := &liquidServiceInfoCmd{}
	cmd := &cobra.Command{
		Use:   "service-info $SERVICE_TYPE",
		Short: "Get information about a liquid",
		Long:  "Get information about a liquid and the resources available within it. This command requires a cloud-admin token.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  liquidServiceInfo.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	liquidServiceInfo.flags.AddToCmd(cmd)

	liquidServiceInfo.Command = cmd
	return liquidServiceInfo
}

func GetLiquidServiceInfo(provider *gophercloud.ProviderClient, opts liquidapi.ClientOpts, ctx context.Context, output bool) (liquid.ServiceInfo, error) {
	liquidClient, err := liquidapi.NewClient(provider, gophercloud.EndpointOpts{}, opts)
	if err != nil {
		return liquid.ServiceInfo{}, util.WrapError(err, "could not instantiate new LIQUID client")
	}
	serviceInfo, err := liquidClient.GetInfo(ctx)
	if err != nil {
		return liquid.ServiceInfo{}, util.WrapError(err, "could not fetch service info from LIQUID api")
	}
	err = liquid.ValidateServiceInfo(serviceInfo)
	if err != nil {
		return liquid.ServiceInfo{}, util.WrapError(err, "received an invalid service info")
	}
	if output {
		err = prettyPrint(serviceInfo)
		if err != nil {
			return liquid.ServiceInfo{}, util.WrapError(err, "could not output service info")
		}
	}
	return serviceInfo, nil
}

func (c *liquidServiceInfoCmd) Run(cmd *cobra.Command, args []string) error {
	var serviceType string
	if len(args) == 1 {
		serviceType = args[0]
	}
	endpoint := c.flags.endpoint
	compare := c.flags.compare
	if compare && (endpoint == "" || serviceType == "") {
		return errors.New("argument $SERVICE_TYPE and flag --endpoint are both required for comparison mode")
	}
	body := c.flags.body
	if body != "" {
		return errors.New("custom request body is not needed when retrieving service info")
	}

	provider, err := authenticate(cmd.Context())
	if err != nil {
		return util.WrapError(err, "could not authenticate with openstack")
	}

	var serviceInfo liquid.ServiceInfo
	if serviceType != "" {
		serviceInfo, err = GetLiquidServiceInfo(provider, liquidapi.ClientOpts{ServiceType: "liquid-" + serviceType}, cmd.Context(), !compare)
		if err != nil {
			return err
		}
	}

	var localServiceInfo liquid.ServiceInfo
	if endpoint != "" {
		localServiceInfo, err = GetLiquidServiceInfo(provider, liquidapi.ClientOpts{EndpointOverride: endpoint}, cmd.Context(), !compare)
		if err != nil {
			return err
		}
	}

	if compare {
		serviceInfo.Version = 0 // Version is usually a timestamp, ignore when comparing responses
		localServiceInfo.Version = 0
		diff := cmp.Diff(serviceInfo, localServiceInfo)
		if diff == "" {
			fmt.Println("ServiceInfo responses are identical")
		} else {
			fmt.Println(diff)
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// LIQUID Report Capacity

type liquidReportCapacityCmd struct {
	*cobra.Command

	flags liquidOperationFlags
}

func newLiquidReportCapacityCmd() *liquidReportCapacityCmd {
	liquidReportCapacity := &liquidReportCapacityCmd{}
	cmd := &cobra.Command{
		Use:     "report-capacity $SERVICE_TYPE",
		Short:   "Get available capacity of a liquid",
		Long:    "Get available capacity across all resources of a liquid. This command requires a cloud-admin token.",
		Args:    cobra.ExactArgs(1),
		PreRunE: authWithLimesAdmin,
		RunE:    liquidReportCapacity.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	liquidReportCapacity.flags.AddToCmd(cmd)

	liquidReportCapacity.Command = cmd
	return liquidReportCapacity
}

func GetLiquidCapacityReport(provider *gophercloud.ProviderClient, opts liquidapi.ClientOpts, ctx context.Context, serviceCapacityRequest *liquid.ServiceCapacityRequest, serviceInfo liquid.ServiceInfo, output bool) (liquid.ServiceCapacityReport, error) {
	liquidClient, err := liquidapi.NewClient(provider, gophercloud.EndpointOpts{}, opts)
	if err != nil {
		return liquid.ServiceCapacityReport{}, util.WrapError(err, "could not instantiate new LIQUID client")
	}
	serviceCapacityReport, err := liquidClient.GetCapacityReport(ctx, *serviceCapacityRequest)
	if err != nil {
		return liquid.ServiceCapacityReport{}, util.WrapError(err, "could not fetch service capacity report from LIQUID api")
	}
	err = liquid.ValidateCapacityReport(serviceCapacityReport, *serviceCapacityRequest, serviceInfo)
	if err != nil {
		return liquid.ServiceCapacityReport{}, util.WrapError(err, "received an invalid service capacity report")
	}
	if output {
		err = prettyPrint(serviceCapacityReport)
		if err != nil {
			return liquid.ServiceCapacityReport{}, util.WrapError(err, "could not output service capacity report")
		}
	}
	return serviceCapacityReport, nil
}

func (c *liquidReportCapacityCmd) Run(cmd *cobra.Command, args []string) error {
	serviceType := args[0]

	endpoint := c.flags.endpoint
	compare := c.flags.compare
	if compare && (endpoint == "" || serviceType == "") {
		return errors.New("argument $SERVICE_TYPE and flag --endpoint are both required for comparison mode")
	}
	body := c.flags.body

	var serviceCapacityRequest *liquid.ServiceCapacityRequest
	if body == "" {
		var res gophercloud.Result
		resp, err := limesAdminClient.Get(cmd.Context(), limesAdminClient.ServiceURL("liquid/service-capacity-request?service_type=liquid-"+serviceType), &res.Body, nil) //nolint:bodyclose
		_, res.Header, res.Err = gophercloud.ParseResponse(resp, err)
		if err != nil {
			return util.WrapError(err, "could not fetch service capacity request body from limes")
		}

		err = res.ExtractInto(&serviceCapacityRequest)
		if err != nil {
			return util.WrapError(err, "could not parse service capacity request body response from limes")
		}
	} else {
		decoder := json.NewDecoder(strings.NewReader(body))
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&serviceCapacityRequest)
		if err != nil {
			return util.WrapError(err, "could not parse custom body parameters into service capacity request")
		}
	}

	provider, err := authenticate(cmd.Context())
	if err != nil {
		return util.WrapError(err, "could not authenticate with openstack")
	}

	var serviceCapacityReport liquid.ServiceCapacityReport
	if compare || endpoint == "" {
		serviceInfo, err := GetLiquidServiceInfo(provider, liquidapi.ClientOpts{ServiceType: "liquid-" + serviceType}, cmd.Context(), false)
		if err != nil {
			return err
		}
		serviceCapacityReport, err = GetLiquidCapacityReport(provider, liquidapi.ClientOpts{ServiceType: "liquid-" + serviceType}, cmd.Context(), serviceCapacityRequest, serviceInfo, !compare)
		if err != nil {
			return err
		}
	}

	var localServiceCapacityReport liquid.ServiceCapacityReport
	if endpoint != "" {
		localServiceInfo, err := GetLiquidServiceInfo(provider, liquidapi.ClientOpts{EndpointOverride: endpoint}, cmd.Context(), false)
		if err != nil {
			return err
		}
		localServiceCapacityReport, err = GetLiquidCapacityReport(provider, liquidapi.ClientOpts{EndpointOverride: endpoint}, cmd.Context(), serviceCapacityRequest, localServiceInfo, !compare)
		if err != nil {
			return err
		}
	}

	if compare {
		serviceCapacityReport.InfoVersion = 0 // InfoVersion is usually a timestamp, ignore when comparing responses
		localServiceCapacityReport.InfoVersion = 0
		diff := cmp.Diff(serviceCapacityReport, localServiceCapacityReport)
		if diff == "" {
			fmt.Println("ServiceCapacityReports are identical")
		} else {
			fmt.Println(diff)
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// LIQUID Report Usage

type liquidReportUsageCmd struct {
	*cobra.Command

	flags liquidOperationFlags
}

func newLiquidReportUsageCmd() *liquidReportUsageCmd {
	liquidReportUsage := &liquidReportUsageCmd{}
	cmd := &cobra.Command{
		Use:     "report-usage $SERVICE_TYPE $PROJECT_UUID",
		Short:   "Get usage data within a project of a liquid",
		Long:    "Get usage data (as well as applicable quotas) within a project across all resources of a liquid. This command requires a cloud-admin token.",
		Args:    cobra.ExactArgs(2),
		PreRunE: authWithLimesAdmin,
		RunE:    liquidReportUsage.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	liquidReportUsage.flags.AddToCmd(cmd)

	liquidReportUsage.Command = cmd
	return liquidReportUsage
}

func GetLiquidUsageReport(provider *gophercloud.ProviderClient, opts liquidapi.ClientOpts, ctx context.Context, projectID string, serviceUsageRequest *liquid.ServiceUsageRequest, serviceInfo liquid.ServiceInfo, output bool) (liquid.ServiceUsageReport, error) {
	liquidClient, err := liquidapi.NewClient(provider, gophercloud.EndpointOpts{}, opts)
	if err != nil {
		return liquid.ServiceUsageReport{}, util.WrapError(err, "could not instantiate new LIQUID client")
	}
	serviceUsageReport, err := liquidClient.GetUsageReport(ctx, projectID, *serviceUsageRequest)
	if err != nil {
		return liquid.ServiceUsageReport{}, util.WrapError(err, "could not fetch service usage report from LIQUID api")
	}
	err = liquid.ValidateUsageReport(serviceUsageReport, *serviceUsageRequest, serviceInfo)
	if err != nil {
		return liquid.ServiceUsageReport{}, util.WrapError(err, "received an invalid service usage report")
	}
	if output {
		err = prettyPrint(serviceUsageReport)
		if err != nil {
			return liquid.ServiceUsageReport{}, util.WrapError(err, "could not output service usage report")
		}
	}
	return serviceUsageReport, nil
}

func (c *liquidReportUsageCmd) Run(cmd *cobra.Command, args []string) error {
	serviceType := args[0]
	projectID := args[1]

	endpoint := c.flags.endpoint
	compare := c.flags.compare
	if compare && (endpoint == "" || serviceType == "") {
		return errors.New("argument $SERVICE_TYPE and flag --endpoint are both required for comparison mode")
	}
	body := c.flags.body

	var serviceUsageRequest *liquid.ServiceUsageRequest
	if body == "" {
		var r gophercloud.Result
		resp, err := limesAdminClient.Get(cmd.Context(), limesAdminClient.ServiceURL(fmt.Sprintf("liquid/service-usage-request?service_type=%s&project_id=%s", serviceType, projectID)), &r.Body, nil) //nolint:bodyclose
		_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
		if err != nil {
			return util.WrapError(err, "could not fetch service usage request body from limes")
		}

		err = r.ExtractInto(&serviceUsageRequest)
		if err != nil {
			return util.WrapError(err, "could not parse service usage request body response from limes")
		}
	} else {
		decoder := json.NewDecoder(strings.NewReader(body))
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&serviceUsageRequest)
		if err != nil {
			return util.WrapError(err, "could not parse custom body parameters into service usage request")
		}
	}

	provider, err := authenticate(cmd.Context())
	if err != nil {
		return util.WrapError(err, "could not authenticate with openstack")
	}

	var serviceUsageReport liquid.ServiceUsageReport
	if compare || endpoint == "" {
		serviceInfo, err := GetLiquidServiceInfo(provider, liquidapi.ClientOpts{ServiceType: "liquid-" + serviceType}, cmd.Context(), false)
		if err != nil {
			return err
		}
		serviceUsageReport, err = GetLiquidUsageReport(provider, liquidapi.ClientOpts{ServiceType: "liquid-" + serviceType}, cmd.Context(), projectID, serviceUsageRequest, serviceInfo, !compare)
		if err != nil {
			return err
		}
	}

	var localServiceusageReport liquid.ServiceUsageReport
	if endpoint != "" {
		localServiceInfo, err := GetLiquidServiceInfo(provider, liquidapi.ClientOpts{EndpointOverride: endpoint}, cmd.Context(), false)
		if err != nil {
			return err
		}
		localServiceusageReport, err = GetLiquidUsageReport(provider, liquidapi.ClientOpts{EndpointOverride: endpoint}, cmd.Context(), projectID, serviceUsageRequest, localServiceInfo, !compare)
		if err != nil {
			return err
		}
	}

	if compare {
		serviceUsageReport.InfoVersion = 0 // InfoVersion is usually a timestamp, ignore when comparing responses
		localServiceusageReport.InfoVersion = 0
		diff := cmp.Diff(serviceUsageReport, localServiceusageReport)
		if diff == "" {
			fmt.Println("ServiceUsageReports are identical")
		} else {
			fmt.Println(diff)
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// LIQUID Set Quota

type liquidSetQuotaCmd struct {
	*cobra.Command

	flags liquidQuotaOperationFlags
}

func newLiquidSetQuotaCmd() *liquidSetQuotaCmd {
	liquidSetQuota := &liquidSetQuotaCmd{}
	cmd := &cobra.Command{
		Use:   "set-quota $SERVICE_TYPE $PROJECT_UUID",
		Short: "Updates quota within a project if a liquid service",
		Long:  "Updates quota within a project across specified resources of a liquid. This command requires a cloud-admin token.",
		Args:  cobra.ExactArgs(2),
		RunE:  liquidSetQuota.Run,
	}

	// Flags
	doNotSortFlags(cmd)
	liquidSetQuota.flags.AddToCmd(cmd)

	liquidSetQuota.Command = cmd
	return liquidSetQuota
}

func (c *liquidSetQuotaCmd) Run(cmd *cobra.Command, args []string) error {
	serviceType := args[0]
	projectID := args[1]

	endpoint := c.flags.endpoint

	if len(c.flags.quotaValues) == 0 {
		return errors.New("flag --quota-values is required")
	}
	var serviceQuotaRequest liquid.ServiceQuotaRequest
	serviceQuotaRequest.Resources = make(map[liquid.ResourceName]liquid.ResourceQuotaRequest)
	resourceQuotaStrings := c.flags.quotaValues
	for _, resourceQuotaString := range resourceQuotaStrings {
		parts := strings.Split(resourceQuotaString, "=")
		if len(parts) != 2 {
			return errors.New("quota values should be formatted like $RESOURCE1=$VALUE1,$RESOURCE2=$VALUE2")
		}
		quota, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return util.WrapError(err, "could not convert quota value to uint64")
		}

		serviceQuotaRequest.Resources[liquid.ResourceName(parts[0])] = liquid.ResourceQuotaRequest{Quota: quota}
	}

	provider, err := authenticate(cmd.Context())
	if err != nil {
		return util.WrapError(err, "could not authenticate with openstack")
	}

	var liquidClient *liquidapi.Client
	if endpoint == "" {
		liquidClient, err = liquidapi.NewClient(provider, gophercloud.EndpointOpts{}, liquidapi.ClientOpts{ServiceType: "liquid-" + serviceType})
	} else {
		liquidClient, err = liquidapi.NewClient(provider, gophercloud.EndpointOpts{}, liquidapi.ClientOpts{EndpointOverride: endpoint})
	}
	if err != nil {
		return util.WrapError(err, "could not instantiate new LIQUID client")
	}

	var serviceInfo liquid.ServiceInfo
	if endpoint == "" {
		serviceInfo, err = GetLiquidServiceInfo(provider, liquidapi.ClientOpts{ServiceType: "liquid-" + serviceType}, cmd.Context(), false)
	} else {
		serviceInfo, err = GetLiquidServiceInfo(provider, liquidapi.ClientOpts{EndpointOverride: endpoint}, cmd.Context(), false)
	}
	if err != nil {
		return err
	}

	// Validate that quotas are provided for all resources
	for resourceName := range serviceInfo.Resources {
		if _, ok := serviceQuotaRequest.Resources[resourceName]; !ok {
			return errors.New("quota missing for resource " + string(resourceName))
		}
	}
	for resourceName := range serviceQuotaRequest.Resources {
		if _, ok := serviceInfo.Resources[resourceName]; !ok {
			return errors.New("quota provided for unknown resource " + string(resourceName))
		}
	}

	err = liquidClient.PutQuota(cmd.Context(), projectID, serviceQuotaRequest)
	if err != nil {
		return util.WrapError(err, "could not set quota values")
	}
	fmt.Println("successfully set quota values")

	return nil
}
