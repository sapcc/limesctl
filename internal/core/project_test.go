package core

import (
	"bytes"
	"encoding/json"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/sapcc/limes"
)

func TestProjectReportRender(t *testing.T) {
	mockJSONBytes, err := getFixtureBytes("project-get-dresden.json")
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
	err = RenderReports(CSVRecordFormatDefault, false, rep).Write(&actual)
	th.AssertNoErr(t, err)

	mockCSVBytes, err := getFixtureBytes("project-get-dresden.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(mockCSVBytes), actual.String())
}

func TestProjectReportsRender(t *testing.T) {
	type listData struct {
		Projects []limes.ProjectReport `json:"projects"`
	}
	domainID := "uuid-for-germany"
	domainName := "germany"

	// List
	mockJSONBytes, err := getFixtureBytes("project-list.json")
	th.AssertNoErr(t, err)
	var data listData
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	var actual bytes.Buffer
	reps := LimesProjectsToReportRenderer(data.Projects, domainID, domainName)
	err = RenderReports(CSVRecordFormatDefault, false, reps...).Write(&actual)
	th.AssertNoErr(t, err)

	mockCSVBytes, err := getFixtureBytes("project-list.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(mockCSVBytes), actual.String())

	// Filtered list with long CSV format and human-readable values
	mockJSONBytes, err = getFixtureBytes("project-list-filtered.json")
	th.AssertNoErr(t, err)
	var filteredData listData
	err = json.Unmarshal(mockJSONBytes, &filteredData)
	th.AssertNoErr(t, err)

	actual.Reset()
	reps = LimesProjectsToReportRenderer(filteredData.Projects, domainID, domainName)
	err = RenderReports(CSVRecordFormatLong, true, reps...).Write(&actual)
	th.AssertNoErr(t, err)

	mockCSVBytes, err = getFixtureBytes("project-list-filtered.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(mockCSVBytes), actual.String())
}
