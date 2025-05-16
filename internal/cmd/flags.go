// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/sapcc/limesctl/v3/internal/core"
)

func doNotSortFlags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	cmd.LocalFlags().SortFlags = false
	cmd.PersistentFlags().SortFlags = false
	cmd.InheritedFlags().SortFlags = false
}

///////////////////////////////////////////////////////////////////////////////
// Limes' API request filter flags.

type commonFilterFlags struct {
	services []string
	areas    []string
}

func (f *commonFilterFlags) AddToCmd(cmd *cobra.Command) {
	cmd.Flags().StringSliceVar(&f.services, "services", nil, "service types (comma separated list)")
	cmd.Flags().StringSliceVar(&f.areas, "areas", nil, "service areas (comma separated list)")
}

// resourceFilterFlags define parameters for Limes API requests that concern
// resources.
type resourceFilterFlags struct {
	commonFilterFlags
	resources []string
}

func (f *resourceFilterFlags) AddToCmd(cmd *cobra.Command) {
	f.commonFilterFlags.AddToCmd(cmd)
	cmd.Flags().StringSliceVar(&f.resources, "resources", nil, "resource names (comma separated list)")
}

// rateFilterFlags define parameters for Limes API requests that concern rates.
type rateFilterFlags struct {
	commonFilterFlags
}

func (f *rateFilterFlags) AddToCmd(cmd *cobra.Command) {
	f.commonFilterFlags.AddToCmd(cmd)
}

///////////////////////////////////////////////////////////////////////////////
// CLI output format flags.

type commonOutputFmtFlags struct {
	format core.OutputFormat
	names  bool
	long   bool
}

func (o *commonOutputFmtFlags) AddToCmd(cmd *cobra.Command) {
	cmd.Flags().VarP(&o.format, "format", "f", "output format: table (default), json, csv")
	cmd.Flags().BoolVar(&o.names, "names", false, "show output with names instead of UUIDs. Not valid for 'json' output format")
	cmd.Flags().BoolVar(&o.long, "long", false, "show detailed output. Not valid for 'json' output format")
}

func (o commonOutputFmtFlags) validate() (*core.OutputOpts, error) {
	// Catch errors.
	if o.long && o.names {
		return nil, errors.New("'--long' and '--names' flags are mutually exclusive, i.e. use one, not both")
	}

	opts := &core.OutputOpts{
		Fmt: o.format,
	}
	switch {
	case o.long:
		opts.CSVRecFmt = core.CSVRecordFormatLong
	case o.names:
		opts.CSVRecFmt = core.CSVRecordFormatNames
	default:
		opts.CSVRecFmt = core.CSVRecordFormatDefault
	}

	return opts, nil
}

// resourceOutputFmtFlags define how the app will print resource data.
type resourceOutputFmtFlags struct {
	commonOutputFmtFlags
	humanize bool
}

func (o *resourceOutputFmtFlags) AddToCmd(cmd *cobra.Command) {
	o.commonOutputFmtFlags.AddToCmd(cmd)
	cmd.Flags().BoolVar(&o.humanize, "humanize", false, "show quota and usage values in an user friendly unit. Not valid for 'json' output format")
}

func (o resourceOutputFmtFlags) validate() (*core.OutputOpts, error) {
	opts, err := o.commonOutputFmtFlags.validate()
	if err != nil {
		return nil, err
	}

	opts.Humanize = o.humanize
	return opts, nil
}

// rateOutputFmtFlags define how the app will print rate limit data.
type rateOutputFmtFlags struct {
	commonOutputFmtFlags
}

func (o *rateOutputFmtFlags) AddToCmd(cmd *cobra.Command) {
	o.commonOutputFmtFlags.AddToCmd(cmd)
}

func (o rateOutputFmtFlags) validate() (*core.OutputOpts, error) {
	return o.commonOutputFmtFlags.validate()
}

// liquidOperationFlags
type liquidOperationFlags struct {
	endpoint string
	compare  bool
	body     string
}

func (l *liquidOperationFlags) AddToCmd(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&l.endpoint, "endpoint", "e", "", "query a liquid running locally")
	cmd.Flags().BoolVarP(&l.compare, "compare", "c", false, "query both the liquid in the cluster and the liquid running locally. Renders the diff of both responses. Requires --endpoint to be set.")
	cmd.Flags().StringVarP(&l.body, "body", "b", "", "use a custom request body generated from provided json structure")
}

// liquidQuotaOperationFlags
type liquidQuotaOperationFlags struct {
	endpoint    string
	quotaValues []string
}

func (l *liquidQuotaOperationFlags) AddToCmd(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&l.endpoint, "endpoint", "e", "", "query a liquid running locally")
	cmd.Flags().StringSliceVarP(&l.quotaValues, "quota-values", "q", nil, "quota values $RESOURCE=$VALUE (comma-separated list)")
}
