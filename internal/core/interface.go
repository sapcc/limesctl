/*******************************************************************************
*
* Copyright 2018-2020 SAP SE
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
	CSVRecordFormatLong CSVRecordFormat = iota + 1
	CSVRecordFormatNames
)

// CSVRecords is exactly that.
type CSVRecords [][]string

// LimesReportRenderer is implemented by data types that can render a Limes
// API report into CSVRecords.
type LimesReportRenderer interface {
	getHeaderRow(csvFmt CSVRecordFormat) []string
	render(csvFmt CSVRecordFormat, humanize bool) CSVRecords
}

// RenderReports renders multiple reports and returns the aggregate CSVData.
func RenderReports(rL []LimesReportRenderer, csvFmt CSVRecordFormat, humanize bool) CSVRecords {
	var recs CSVRecords
	recs = append(recs, rL[0].getHeaderRow(csvFmt))
	for _, r := range rL {
		recs = append(recs, r.render(csvFmt, humanize)...)
	}
	return recs
}
