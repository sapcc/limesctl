/*******************************************************************************
*
* Copyright 2018 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package cli

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/sapcc/limes"
)

var mockRawCapacities = []string{
	"compute/cores=10",
	"compute/ram=20MiB",
	"object-store/capacity=30B:I got 99 problems, but a cluster ain't one.",
}

var mockRawQuotas = []string{
	"compute/cores=10",
	"compute/ram=20MiB",
	"object-store/capacity=30B:this comment should be ignored by the parser.",
}

// TestParseCapacities tests if the service capacities given at the
// command line are being correctly parsed.
func TestParseCapacities(t *testing.T) {
	q := &Quotas{}
	for _, v := range mockRawCapacities {
		err := q.Set(v)
		th.AssertNoErr(t, err)
	}
	actual := makeServiceCapacities(q)

	unitB := limes.UnitBytes
	unitMiB := limes.UnitMebibytes
	unitNone := limes.UnitNone
	expected := []limes.ServiceCapacityRequest{
		{Type: "compute", Resources: []limes.ResourceCapacityRequest{
			{
				Name:     "cores",
				Capacity: 10,
				Unit:     &unitNone,
			},
			{
				Name:     "ram",
				Capacity: 20,
				Unit:     &unitMiB,
			},
		}},
		{Type: "object-store", Resources: []limes.ResourceCapacityRequest{
			{
				Name:     "capacity",
				Capacity: 30,
				Unit:     &unitB,
				Comment:  "I got 99 problems, but a cluster ain't one.",
			},
		}},
	}
	th.AssertDeepEquals(t, expected, actual)
}

// TestParseQuotas tests if the service quotas given at the command line
// are being correctly parsed.
func TestParseQuotas(t *testing.T) {
	q := &Quotas{}
	for _, v := range mockRawQuotas {
		err := q.Set(v)
		th.AssertNoErr(t, err)
	}
	actual := makeServiceQuotas(q)

	expected := limes.QuotaRequest{
		"compute": limes.ServiceQuotaRequest{
			"cores": limes.ValueWithUnit{10, limes.UnitNone},
			"ram":   limes.ValueWithUnit{20, limes.UnitMebibytes},
		},
		"object-store": limes.ServiceQuotaRequest{
			"capacity": limes.ValueWithUnit{30, limes.UnitBytes},
		},
	}
	th.AssertDeepEquals(t, expected, actual)
}

func TestRenderClusterCSV(t *testing.T) {
	var actual bytes.Buffer

	// get
	c, err := makeMockCluster("./fixtures/cluster-get.json")
	th.AssertNoErr(t, err)

	expected, err := ioutil.ReadFile("./fixtures/cluster-get.csv")
	th.AssertNoErr(t, err)

	c.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())

	// filtered get
	c, err = makeMockCluster("./fixtures/cluster-get-filtered.json")
	th.AssertNoErr(t, err)
	c.Output = Output{
		HumanReadable: true,
		Long:          true,
	}
	c.Filter = Filter{
		Area:     "shared",
		Service:  "shared",
		Resource: "capacity",
	}

	expected, err = ioutil.ReadFile("./fixtures/cluster-get-filtered.csv")
	th.AssertNoErr(t, err)

	actual.Reset()
	c.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())

	// list
	c, err = makeMockCluster("./fixtures/cluster-list.json")
	th.AssertNoErr(t, err)
	c.IsList = true

	expected, err = ioutil.ReadFile("./fixtures/cluster-list.csv")
	th.AssertNoErr(t, err)

	actual.Reset()
	c.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())

	// filtered list
	c, err = makeMockCluster("./fixtures/cluster-list-filtered.json")
	th.AssertNoErr(t, err)
	c.IsList = true
	c.Output = Output{
		Long: true,
	}
	c.Filter = Filter{
		Area:     "shared",
		Service:  "shared",
		Resource: "capacity",
	}

	expected, err = ioutil.ReadFile("./fixtures/cluster-list-filtered.csv")
	th.AssertNoErr(t, err)

	actual.Reset()
	c.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())
}

func TestRenderDomainCSV(t *testing.T) {
	var actual bytes.Buffer

	// get
	d, err := makeMockDomain("./fixtures/domain-get.json")
	th.AssertNoErr(t, err)

	expected, err := ioutil.ReadFile("./fixtures/domain-get.csv")
	th.AssertNoErr(t, err)

	d.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())

	// filtered get
	d, err = makeMockDomain("./fixtures/domain-get-filtered.json")
	th.AssertNoErr(t, err)
	d.Output = Output{
		HumanReadable: true,
		Names:         true,
	}
	d.Filter = Filter{
		Area:     "shared",
		Service:  "shared",
		Resource: "capacity",
	}

	expected, err = ioutil.ReadFile("./fixtures/domain-get-filtered.csv")
	th.AssertNoErr(t, err)

	actual.Reset()
	d.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())

	// list
	d, err = makeMockDomain("./fixtures/domain-list.json")
	th.AssertNoErr(t, err)
	d.IsList = true

	expected, err = ioutil.ReadFile("./fixtures/domain-list.csv")
	th.AssertNoErr(t, err)

	actual.Reset()
	d.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())

	// filtered list
	d, err = makeMockDomain("./fixtures/domain-list-filtered.json")
	th.AssertNoErr(t, err)
	d.IsList = true
	d.Output = Output{
		Long: true,
	}
	d.Filter = Filter{
		Area:     "shared",
		Service:  "shared",
		Resource: "things",
	}

	expected, err = ioutil.ReadFile("./fixtures/domain-list-filtered.csv")
	th.AssertNoErr(t, err)

	actual.Reset()
	d.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())
}

func TestRenderProjectCSV(t *testing.T) {
	var actual bytes.Buffer

	// get
	p, err := makeMockProject("./fixtures/project-get.json")
	th.AssertNoErr(t, err)

	expected, err := ioutil.ReadFile("./fixtures/project-get.csv")
	th.AssertNoErr(t, err)

	p.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())

	// filtered get
	p, err = makeMockProject("./fixtures/project-get-filtered.json")
	th.AssertNoErr(t, err)
	p.Output = Output{
		HumanReadable: true,
		Names:         true,
	}
	p.Filter = Filter{
		Area:     "shared",
		Service:  "shared",
		Resource: "capacity",
	}

	expected, err = ioutil.ReadFile("./fixtures/project-get-filtered.csv")
	th.AssertNoErr(t, err)

	actual.Reset()
	p.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())

	// list
	p, err = makeMockProject("./fixtures/project-list.json")
	th.AssertNoErr(t, err)
	p.IsList = true

	expected, err = ioutil.ReadFile("./fixtures/project-list.csv")
	th.AssertNoErr(t, err)

	actual.Reset()
	p.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())

	// filtered list
	p, err = makeMockProject("./fixtures/project-list-filtered.json")
	th.AssertNoErr(t, err)
	p.IsList = true
	p.Output = Output{
		Long: true,
	}
	p.Filter = Filter{
		Area:     "shared",
		Service:  "shared",
		Resource: "things",
	}

	expected, err = ioutil.ReadFile("./fixtures/project-list-filtered.csv")
	th.AssertNoErr(t, err)

	actual.Reset()
	p.renderCSV().write(&actual)
	th.AssertEquals(t, string(expected), actual.String())
}

// makeMockCluster is a helper function that uses a JSON file to create a mock
// Cluster for testing and assigns the unserialized JSON as its Result.Body.
func makeMockCluster(pathToJSON string) (*Cluster, error) {
	c := new(Cluster)

	b, err := ioutil.ReadFile(pathToJSON)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &c.Result.Body)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// makeMockDomain is a helper function that uses a JSON file to create a mock
// Domain for testing and assigns the unserialized JSON as its Result.Body.
func makeMockDomain(pathToJSON string) (*Domain, error) {
	d := new(Domain)

	b, err := ioutil.ReadFile(pathToJSON)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &d.Result.Body)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// makeMockProject is a helper function that uses a JSON file to create a mock
// Project for testing and assigns the unserialized JSON as its Result.Body.
func makeMockProject(pathToJSON string) (*Project, error) {
	p := new(Project)
	p.DomainID = "uuid-for-germany"
	p.DomainName = "germany"

	b, err := ioutil.ReadFile(pathToJSON)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &p.Result.Body)
	if err != nil {
		return nil, err
	}

	return p, nil
}
