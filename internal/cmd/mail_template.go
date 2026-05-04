// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/spf13/cobra"

	"github.com/sapcc/limesctl/v3/internal/util"
)

// mailTemplates holds the rendered mail templates from the API.
type mailTemplates struct {
	ConfirmedCommitments   string `json:"confirmed_commitments"`
	ExpiringCommitments    string `json:"expiring_commitments"`
	TransferredCommitments string `json:"transferred_commitments"`
}

func newMailTemplateCmd() *cobra.Command {
	var fetchedMailTemplates *mailTemplates

	// subcommands can access the templates.
	// called after PersistentPreRunE has populated the data.
	getTemplates := func() *mailTemplates {
		return fetchedMailTemplates
	}

	cmd := &cobra.Command{
		Use:   "mail-template",
		Short: "Display configured mail templates (cloud-admin only)",
		Long:  "Display configured mail templates. Can be used to check the validity of the configured templates.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := authWithLimesAdmin(cmd, args); err != nil {
				return err
			}

			var res gophercloud.Result
			resp, err := limesAdminClient.Get(cmd.Context(), limesAdminClient.ServiceURL("mail/render"), &res.Body, nil) //nolint:bodyclose
			_, _, res.Err = gophercloud.ParseResponse(resp, err)
			if res.Err != nil {
				return util.WrapError(res.Err, "could not fetch mail templates from limes")
			}

			fetchedMailTemplates = &mailTemplates{}
			if err := res.ExtractInto(fetchedMailTemplates); err != nil {
				return util.WrapError(err, "could not extract mail templates")
			}

			return nil
		},
	}

	doNotSortFlags(cmd)

	// Add subcommands
	cmd.AddCommand(newAllTemplatesCmd(getTemplates))
	cmd.AddCommand(newConfirmedCommitmentsCmd(getTemplates))
	cmd.AddCommand(newExpiringCommitmentsCmd(getTemplates))
	cmd.AddCommand(newTransferredCommitmentsCmd(getTemplates))

	return cmd
}

func newAllTemplatesCmd(getTemplates func() *mailTemplates) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Display all mail templates",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			templates := getTemplates()

			fmt.Println("=== Confirmed Commitments ===")
			fmt.Println(templates.ConfirmedCommitments)

			fmt.Println("=== Expiring Commitments ===")
			fmt.Println(templates.ExpiringCommitments)

			fmt.Println("=== Transferred Commitments ===")
			fmt.Println(templates.TransferredCommitments)

			return nil
		},
	}

	doNotSortFlags(cmd)
	return cmd
}

func newConfirmedCommitmentsCmd(getTemplates func() *mailTemplates) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirmed",
		Short: "Display the confirmed commitments mail template",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return outputTemplate(getTemplates().ConfirmedCommitments)
		},
	}

	doNotSortFlags(cmd)
	return cmd
}

func newExpiringCommitmentsCmd(getTemplates func() *mailTemplates) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expiring",
		Short: "Display the expiring commitments mail template",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return outputTemplate(getTemplates().ExpiringCommitments)
		},
	}

	doNotSortFlags(cmd)
	return cmd
}

func newTransferredCommitmentsCmd(getTemplates func() *mailTemplates) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transferred",
		Short: "Display the transferred commitments mail template",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return outputTemplate(getTemplates().TransferredCommitments)
		},
	}

	doNotSortFlags(cmd)
	return cmd
}

// outputTemplate handles the output of a single template.
func outputTemplate(content string) error {
	fmt.Println(content)
	return nil
}
