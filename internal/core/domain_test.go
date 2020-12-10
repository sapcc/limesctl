package core

import (
	"bytes"
	"encoding/json"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/sapcc/limes"
)

func TestDomainReportRender(t *testing.T) {
	mockJSONBytes, err := getFixtureBytes("domain-get-germany.json")
	th.AssertNoErr(t, err)
	var data struct {
		Domain limes.DomainReport `json:"domain"`
	}
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	var actual bytes.Buffer
	rep := DomainReport{&data.Domain}
	err = RenderReports(CSVRecordFormatDefault, false, rep).Write(&actual)
	th.AssertNoErr(t, err)

	mockCSVBytes, err := getFixtureBytes("domain-get-germany.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(mockCSVBytes), actual.String())
}

func TestDomainReportsRender(t *testing.T) {
	type listData struct {
		Domains []limes.DomainReport `json:"domains"`
	}

	// List
	mockJSONBytes, err := getFixtureBytes("domain-list.json")
	th.AssertNoErr(t, err)
	var data listData
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	var actual bytes.Buffer
	reps := LimesDomainsToReportRenderer(data.Domains)
	err = RenderReports(CSVRecordFormatDefault, false, reps...).Write(&actual)
	th.AssertNoErr(t, err)

	mockCSVBytes, err := getFixtureBytes("domain-list.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(mockCSVBytes), actual.String())

	// Filtered list with long CSV format and human-readable values
	mockJSONBytes, err = getFixtureBytes("domain-list-filtered.json")
	th.AssertNoErr(t, err)
	var filteredData listData
	err = json.Unmarshal(mockJSONBytes, &filteredData)
	th.AssertNoErr(t, err)

	actual.Reset()
	reps = LimesDomainsToReportRenderer(filteredData.Domains)
	err = RenderReports(CSVRecordFormatLong, true, reps...).Write(&actual)
	th.AssertNoErr(t, err)

	mockCSVBytes, err = getFixtureBytes("domain-list-filtered.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(mockCSVBytes), actual.String())
}
