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
	"io"
	"io/ioutil"
	"os"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/sapcc/limes/pkg/api"
	"github.com/sapcc/limes/pkg/limes"
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
	expected := []api.ServiceCapacities{
		{Type: "compute", Resources: []api.ResourceCapacity{
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
		{Type: "object-store", Resources: []api.ResourceCapacity{
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

	expected := api.ServiceQuotas{
		"compute": api.ResourceQuotas{
			"cores": limes.ValueWithUnit{10, limes.UnitNone},
			"ram":   limes.ValueWithUnit{20, limes.UnitMebibytes},
		},
		"object-store": api.ResourceQuotas{
			"capacity": limes.ValueWithUnit{30, limes.UnitBytes},
		},
	}
	th.AssertDeepEquals(t, expected, actual)
}

func TestRenderClusterCSV(t *testing.T) {
	// get
	c, err := makeMockCluster("./fixtures/cluster-get.json")
	th.AssertNoErr(t, err)

	actual, err := captureOutput(func() { c.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err := ioutil.ReadFile("./fixtures/cluster-get.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)

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

	actual, err = captureOutput(func() { c.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err = ioutil.ReadFile("./fixtures/cluster-get-filtered.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)

	// list
	c, err = makeMockCluster("./fixtures/cluster-list.json")
	th.AssertNoErr(t, err)
	c.IsList = true

	actual, err = captureOutput(func() { c.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err = ioutil.ReadFile("./fixtures/cluster-list.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)

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

	actual, err = captureOutput(func() { c.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err = ioutil.ReadFile("./fixtures/cluster-list-filtered.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)
}

func TestRenderDomainCSV(t *testing.T) {
	// get
	d, err := makeMockDomain("./fixtures/domain-get.json")
	th.AssertNoErr(t, err)

	actual, err := captureOutput(func() { d.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err := ioutil.ReadFile("./fixtures/domain-get.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)

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

	actual, err = captureOutput(func() { d.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err = ioutil.ReadFile("./fixtures/domain-get-filtered.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)

	// list
	d, err = makeMockDomain("./fixtures/domain-list.json")
	th.AssertNoErr(t, err)
	d.IsList = true

	actual, err = captureOutput(func() { d.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err = ioutil.ReadFile("./fixtures/domain-list.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)

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

	actual, err = captureOutput(func() { d.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err = ioutil.ReadFile("./fixtures/domain-list-filtered.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)
}

func TestRenderProjectCSV(t *testing.T) {
	// get
	p, err := makeMockProject("./fixtures/project-get.json")
	th.AssertNoErr(t, err)

	actual, err := captureOutput(func() { p.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err := ioutil.ReadFile("./fixtures/project-get.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)

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

	actual, err = captureOutput(func() { p.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err = ioutil.ReadFile("./fixtures/project-get-filtered.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)

	// list
	p, err = makeMockProject("./fixtures/project-list.json")
	th.AssertNoErr(t, err)
	p.IsList = true

	actual, err = captureOutput(func() { p.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err = ioutil.ReadFile("./fixtures/project-list.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)

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

	actual, err = captureOutput(func() { p.renderCSV().writeCSV() })
	th.AssertNoErr(t, err)

	expected, err = ioutil.ReadFile("./fixtures/project-list-filtered.csv")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), actual)
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

// captureOutput is a helper function that returns the output
// of os.Stdout as a string.
func captureOutput(do func()) (string, error) {
	// make a copy of Stdout
	old := os.Stdout
	// pipe Stdout...
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	do()
	outC := make(chan bytes.Buffer)
	// make separate goroutine so that Copy
	// doesn't block indefinitely
	go func() {
		var b bytes.Buffer
		io.Copy(&b, r)
		outC <- b
	}()

	// restore Stdout
	w.Close()
	os.Stdout = old
	out := <-outC

	return out.String(), nil
}
