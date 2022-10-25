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
