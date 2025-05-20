// SPDX-FileCopyrightText: 2018 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/sapcc/limesctl/v3/internal/util"
)

// OutputFormat that the app can print data in.
type OutputFormat string

// Different types of OutputFormat.
const (
	OutputFormatTable OutputFormat = "table"
	OutputFormatCSV   OutputFormat = "csv"
	OutputFormatJSON  OutputFormat = "json"
)

// String implements the pflag.Value interface.
func (f *OutputFormat) String() string {
	return string(*f)
}

// Set implements the pflag.Value interface.
func (f *OutputFormat) Set(v string) error {
	switch vf := OutputFormat(v); vf {
	case OutputFormatTable, OutputFormatCSV, OutputFormatJSON:
		*f = vf
		return nil
	default:
		return fmt.Errorf("must be one of [%s, %s, %s], got %s",
			OutputFormatTable, OutputFormatCSV, OutputFormatJSON, v)
	}
}

// Type implements the pflag.Value interface.
func (f *OutputFormat) Type() string {
	return "string"
}

// CSVRecordFormat type defines the style of CSV records.
type CSVRecordFormat int

// Different types of CSVRecordFormat.
const (
	CSVRecordFormatDefault CSVRecordFormat = iota
	CSVRecordFormatLong
	CSVRecordFormatNames
)

// CSVRecords is exactly that.
type CSVRecords [][]string

// Write writes CSVRecords to w.
//
// Note: the method takes an io.Writer because it is used in unit tests.
func (d CSVRecords) Write(w io.Writer) error {
	csvW := csv.NewWriter(w)
	csvW.Comma = rune(';') // Use semicolon as delimiter
	if err := csvW.WriteAll(d); err != nil {
		return util.WrapError(err, "could not write CSV data")
	}
	return nil
}

// WriteAsTable writes CSVRecords to os.Stdout in table format.
func (d CSVRecords) WriteAsTable() {
	t := tablewriter.NewWriter(os.Stdout)
	t.Header(d[0])
	t.Bulk(d[1:]) //nolint:errcheck
	t.Render()    //nolint:errcheck
}

type OutputOpts struct {
	Fmt       OutputFormat
	CSVRecFmt CSVRecordFormat
	Humanize  bool
}

// LimesReportRenderer is implemented by data types that can render a Limes
// API report into CSVRecords.
type LimesReportRenderer interface {
	getHeaderRow(opts *OutputOpts) []string
	render(opts *OutputOpts) CSVRecords
}

// RenderReports renders multiple reports and returns the aggregate CSVRecords.
//
// Note: if multiple LimesReportRenderer are given then they must have the same underlying
// type.
func RenderReports(opts *OutputOpts, rL ...LimesReportRenderer) CSVRecords {
	var recs CSVRecords
	if len(rL) > 0 {
		recs = append(recs, rL[0].getHeaderRow(opts))
		for _, r := range rL {
			recs = append(recs, r.render(opts)...)
		}
	}
	return recs
}
