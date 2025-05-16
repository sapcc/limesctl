// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"bytes"
	"encoding/json"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
)

func TestProjectResourcesSingleReportRender(t *testing.T) {
	mockJSONBytes, err := fixtureBytes("project-get-dresden.json")
	th.AssertNoErr(t, err)
	var data struct {
		Project limesresources.ProjectReport `json:"project"`
	}
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	var actual bytes.Buffer
	rep := ProjectResourcesReport{
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

func TestProjectResourcesMultipleReportsRender(t *testing.T) {
	type listData struct {
		Projects []limesresources.ProjectReport `json:"projects"`
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
	reps := LimesProjectResourcesToReportRenderer(data.Projects, domainID, domainName, false)
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
	reps = LimesProjectResourcesToReportRenderer(filteredData.Projects, domainID, domainName, false)
	err = RenderReports(opts, reps...).Write(&actual)
	th.AssertNoErr(t, err)
	assertEquals(t, "project-list-filtered.csv", actual.Bytes())
}
