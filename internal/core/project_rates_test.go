// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"bytes"
	"encoding/json"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	limesrates "github.com/sapcc/go-api-declarations/limes/rates"
)

func TestProjectRatesSingleReportRender(t *testing.T) {
	mockJSONBytes, err := fixtureBytes("project-get-berlin-only-rates.json")
	th.AssertNoErr(t, err)
	var data struct {
		Project limesrates.ProjectReport `json:"project"`
	}
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	var actual bytes.Buffer
	rep := ProjectRatesReport{
		ProjectReport: &data.Project,
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
