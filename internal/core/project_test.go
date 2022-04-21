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

package core

import (
	"bytes"
	"encoding/json"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/sapcc/go-api-declarations/limes"
)

func TestProjectReportRender(t *testing.T) {
	mockJSONBytes, err := fixtureBytes("project-get-dresden.json")
	th.AssertNoErr(t, err)
	var data struct {
		Project limes.ProjectReport `json:"project"`
	}
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	var actual bytes.Buffer
	rep := ProjectReport{
		ProjectReport: &data.Project,
		DomainID:      "uuid-for-germany",
		DomainName:    "germany",
	}

	opts := &OutputOpts{
		CSVRecFmt: CSVRecordFormatDefault,
		Humanize:  false,
	}
	err = RenderReports(opts, rep).Write(&actual)
	th.AssertNoErr(t, err)
	assertEquals(t, "project-get-dresden.csv", actual.Bytes())
}

func TestProjectRatesReportRender(t *testing.T) {
	mockJSONBytes, err := fixtureBytes("project-get-berlin-only-rates.json")
	th.AssertNoErr(t, err)
	var data struct {
		Project limes.ProjectReport `json:"project"`
	}
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	var actual bytes.Buffer
	rep := ProjectReport{
		ProjectReport: &data.Project,
		HasRatesOnly:  true,
		DomainID:      "uuid-for-germany",
		DomainName:    "germany",
	}

	opts := &OutputOpts{
		CSVRecFmt: CSVRecordFormatLong,
		Humanize:  false,
	}
	err = RenderReports(opts, rep).Write(&actual)
	th.AssertNoErr(t, err)
	assertEquals(t, "project-get-berlin-only-rates.csv", actual.Bytes())
}

func TestProjectReportsRender(t *testing.T) {
	type listData struct {
		Projects []limes.ProjectReport `json:"projects"`
	}
	domainID := "uuid-for-germany"
	domainName := "germany"

	// List
	mockJSONBytes, err := fixtureBytes("project-list.json")
	th.AssertNoErr(t, err)
	var data listData
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	opts := &OutputOpts{
		CSVRecFmt: CSVRecordFormatDefault,
		Humanize:  false,
	}
	var actual bytes.Buffer
	reps := LimesProjectsToReportRenderer(data.Projects, domainID, domainName, false)
	err = RenderReports(opts, reps...).Write(&actual)
	th.AssertNoErr(t, err)
	assertEquals(t, "project-list.csv", actual.Bytes())

	// Filtered list with long CSV format and human-readable values
	mockJSONBytes, err = fixtureBytes("project-list-filtered.json")
	th.AssertNoErr(t, err)
	var filteredData listData
	err = json.Unmarshal(mockJSONBytes, &filteredData)
	th.AssertNoErr(t, err)

	opts.CSVRecFmt = CSVRecordFormatLong
	opts.Humanize = true
	actual.Reset()
	reps = LimesProjectsToReportRenderer(filteredData.Projects, domainID, domainName, false)
	err = RenderReports(opts, reps...).Write(&actual)
	th.AssertNoErr(t, err)
	assertEquals(t, "project-list-filtered.csv", actual.Bytes())
}
