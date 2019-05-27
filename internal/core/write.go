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
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/sapcc/limesctl/internal/errors"
)

// writeJSON is a helper function that writes the JSON data to os.Stdout.
func (data *jsonData) write(writer io.Writer) {
	fmt.Fprintln(writer, string(*data))
}

// Write writes the CSV data to w.
func (data *csvData) write(w io.Writer) {
	csvW := csv.NewWriter(w)
	csvW.Comma = rune(';') // use semicolon as delimiter

	err := csvW.WriteAll(*data)
	errors.Handle(err, "could not write CSV data")
}

// writeTable is a helper function that writes the CSV data to os.Stdout in an ASCII table format.
func (data *csvData) writeTable(writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader((*data)[0])

	for _, v := range (*data)[1:] {
		table.Append(v)
	}
	table.Render()
}
