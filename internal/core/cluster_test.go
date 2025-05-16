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

func TestClusterResourcesSingleReportRender(t *testing.T) {
	mockJSONBytes, err := fixtureBytes("cluster-get-west.json")
	th.AssertNoErr(t, err)
	var data struct {
		Cluster limesresources.ClusterReport `json:"cluster"`
	}
	err = json.Unmarshal(mockJSONBytes, &data)
	th.AssertNoErr(t, err)

	opts := &OutputOpts{
		CSVRecFmt: CSVRecordFormatDefault,
		Humanize:  false,
	}
	var actual bytes.Buffer
	rep := ClusterReport{&data.Cluster}
	err = RenderReports(opts, rep).Write(&actual)
	th.AssertNoErr(t, err)
	assertEquals(t, "cluster-get-west.csv", actual.Bytes())
}
