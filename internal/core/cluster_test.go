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
	"io/ioutil"
	"path/filepath"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/sapcc/limes"
)

func getFixtureBytes(fileName string) ([]byte, error) {
	path := filepath.Join("fixtures", fileName)
	return ioutil.ReadFile(path)
}

//nolint:dupl
func TestClusterReportRender(t *testing.T) {
	mockJSONBytes, err := getFixtureBytes("cluster-get-west.json")
	th.AssertNoErr(t, err)
	var data struct {
		Cluster limes.ClusterReport `json:"cluster"`
	}
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	var actual bytes.Buffer
	rep := ClusterReport{&data.Cluster}
	err = RenderReports(CSVRecordFormatDefault, false, rep).Write(&actual)
	th.AssertNoErr(t, err)

	mockCSVBytes, err := getFixtureBytes("cluster-get-west.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(mockCSVBytes), actual.String())
}

func TestClusterReportsRender(t *testing.T) {
	type listData struct {
		CurrentCluster string                `json:"current_cluster"`
		Clusters       []limes.ClusterReport `json:"clusters"`
	}

	// List
	mockJSONBytes, err := getFixtureBytes("cluster-list.json")
	th.AssertNoErr(t, err)
	var data listData
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	var actual bytes.Buffer
	reps := LimesClustersToReportRenderer(data.Clusters)
	err = RenderReports(CSVRecordFormatDefault, false, reps...).Write(&actual)
	th.AssertNoErr(t, err)

	mockCSVBytes, err := getFixtureBytes("cluster-list.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(mockCSVBytes), actual.String())

	// Filtered list with long CSV format and human-readable values
	mockJSONBytes, err = getFixtureBytes("cluster-list-filtered.json")
	th.AssertNoErr(t, err)
	var filteredData listData
	err = json.Unmarshal(mockJSONBytes, &filteredData)
	th.AssertNoErr(t, err)

	actual.Reset()
	reps = LimesClustersToReportRenderer(filteredData.Clusters)
	err = RenderReports(CSVRecordFormatLong, true, reps...).Write(&actual)
	th.AssertNoErr(t, err)

	mockCSVBytes, err = getFixtureBytes("cluster-list-filtered.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(mockCSVBytes), actual.String())
}
