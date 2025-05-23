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

func TestDomainResourcesSingleReportRender(t *testing.T) {
	mockJSONBytes, err := fixtureBytes("domain-get-germany.json")
	th.AssertNoErr(t, err)
	var data struct {
		Domain limesresources.DomainReport `json:"domain"`
	}
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	opts := &OutputOpts{
		CSVRecFmt: CSVRecordFormatDefault,
		Humanize:  false,
	}
	var actual bytes.Buffer
	rep := DomainReport{&data.Domain}
	err = RenderReports(opts, rep).Write(&actual)
	th.AssertNoErr(t, err)
	assertEquals(t, "domain-get-germany.csv", actual.Bytes())
}

func TestDomainResourcesMultipleReportsRender(t *testing.T) {
	type listData struct {
		Domains []limesresources.DomainReport `json:"domains"`
	}

	// List
	mockJSONBytes, err := fixtureBytes("domain-list.json")
	th.AssertNoErr(t, err)
	var data listData
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	opts := &OutputOpts{
		CSVRecFmt: CSVRecordFormatDefault,
		Humanize:  false,
	}
	var actual bytes.Buffer
	reps := LimesDomainsToReportRenderer(data.Domains)
	err = RenderReports(opts, reps...).Write(&actual)
	th.AssertNoErr(t, err)
	assertEquals(t, "domain-list.csv", actual.Bytes())

	// Filtered list with long CSV format and human-readable values
	mockJSONBytes, err = fixtureBytes("domain-list-filtered.json")
	th.AssertNoErr(t, err)
	var filteredData listData
	err = json.Unmarshal(mockJSONBytes, &filteredData)
	th.AssertNoErr(t, err)

	opts.CSVRecFmt = CSVRecordFormatLong
	opts.Humanize = true
	actual.Reset()
	reps = LimesDomainsToReportRenderer(filteredData.Domains)
	err = RenderReports(opts, reps...).Write(&actual)
	th.AssertNoErr(t, err)
	assertEquals(t, "domain-list-filtered.csv", actual.Bytes())
}
