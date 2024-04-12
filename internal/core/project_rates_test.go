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
