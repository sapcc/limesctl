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

package cli

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sapcc/gophercloud-limes/resources/v1/clusters"
	"github.com/sapcc/gophercloud-limes/resources/v1/domains"
	"github.com/sapcc/gophercloud-limes/resources/v1/projects"
)

// Cluster, initially, contains arguments and flags data for a cluster subcommand (e.g. $ limesctl cluster list).
// As different methods are called on it, the fields within the structure are updated accordingly.
// Call its appropriate method to get/list/update a Cluster.
type Cluster struct {
	ID     string
	Opts   Options
	Result clusters.CommonResult
	IsList bool
}

// Domain, initially, contains arguments and flags data for a domain subcommand (e.g. $ limesctl domain list).
// As different methods are called on it, the fields within the structure are updated accordingly.
// Call its appropriate method to get/list/update a Domain.
type Domain struct {
	ID     string
	Name   string
	Opts   Options
	Result domains.CommonResult
	IsList bool
}

// Project, initially, contains arguments and flags data for a project subcommand (e.g. $ limesctl project show my-project).
// As different methods are called on it, the fields within the structure are updated accordingly.
// Call its appropriate method to get/list/update a Project.
type Project struct {
	ID         string
	Name       string
	DomainID   string
	DomainName string
	Opts       Options
	Result     projects.CommonResult
	IsList     bool
}

// Options contains different options that affect the output of an list operation.
type Options struct {
	Names    bool
	Long     bool
	Area     string
	Service  string
	Resource string
}

// GetTask is the interface type that abstracts a get operation.
type GetTask interface {
	get()
	writeJSON()
	renderCSV() *csvData
}

// RunGetTask is the function that operates on a GetTask and shows the output in the respective
// format that is specified at the command line.
func RunGetTask(t GetTask, outputFmt string) {
	t.get()
	switch outputFmt {
	case "json":
		t.writeJSON()
	case "csv":
		t.renderCSV().writeCSV()
	default:
		t.renderCSV().writeTable()
	}
}

// ListTask is the interface type that abstracts a list operation.
type ListTask interface {
	list()
	writeJSON()
	renderCSV() *csvData
}

// RunListTask is the function that operates on a ListTask and shows the output in the respective
// format that is specified at the command line.
func RunListTask(t ListTask, outputFmt string) {
	t.list()
	switch outputFmt {
	case "json":
		t.writeJSON()
	case "csv":
		t.renderCSV().writeCSV()
	default:
		t.renderCSV().writeTable()
	}
}

// SetTask is the interface type that abstracts a put operation.
type SetTask interface {
	set(interface{})
	writeJSON()
	renderCSV() *csvData
}

// RunSetTask is the function that operates on a SetTask and shows the output in the respective
// format that is specified at the command line.
func RunSetTask(t SetTask, q interface{}, outputFmt string) {
	t.set(q)
	switch outputFmt {
	case "json":
		t.writeJSON()
	case "csv":
		t.renderCSV().writeCSV()
	default:
		t.renderCSV().writeTable()
	}
}

// writeCSV is a helper function that writes the CSV data to os.Stdout.
func (data *csvData) writeCSV() {
	for _, record := range *data {
		var str string
		for i, v := range record {
			// replace empty values with proper CSV null values ("")
			if v == "" {
				v = "\"\""
			}
			// add delimiter for all values except the last one
			if i != (len(record) - 1) {
				v = fmt.Sprintf("%v,", v)
			}
			str += v
		}
		fmt.Println(str)
	}
}

// writeTable is a helper function that writes the CSV data to os.Stdout in an ASCII table format.
func (data *csvData) writeTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader((*data)[0])

	for _, v := range (*data)[1:] {
		table.Append(v)
	}
	table.Render()
}
