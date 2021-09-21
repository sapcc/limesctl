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

package core

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

// OutputFormat that the app can print data in.
type OutputFormat string

// Different types of OutputFormat.
const (
	OutputFormatTable OutputFormat = "table"
	OutputFormatCSV   OutputFormat = "csv"
	OutputFormatJSON  OutputFormat = "json"
)

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
		return errors.Wrap(err, "could not write CSV data")
	}
	return nil
}

// WriteAsTable writes CSVRecords to os.Stdout in table format.
func (d CSVRecords) WriteAsTable() {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader(d[0])
	t.AppendBulk(d[1:])
	t.Render()
}

// LimesReportRenderer is implemented by data types that can render a Limes
// API report into CSVRecords.
type LimesReportRenderer interface {
	getHeaderRow(csvFmt CSVRecordFormat) []string
	render(csvFmt CSVRecordFormat, humanize bool) CSVRecords
}

// RenderReports renders multiple reports and returns the aggregate CSVRecords.
//
// Note: this function expects all LimesReportRenderer to have the same
// underlying type.
func RenderReports(csvFmt CSVRecordFormat, humanize bool, rL ...LimesReportRenderer) CSVRecords {
	var recs CSVRecords
	if len(rL) > 0 {
		recs = append(recs, rL[0].getHeaderRow(csvFmt))
		for _, r := range rL {
			recs = append(recs, r.render(csvFmt, humanize)...)
		}
	}
	return recs
}
