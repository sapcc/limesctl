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

package core

import (
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/sapcc/limes"
)

// TestQuotaValueRx tests if the regular expression used to parse the quota
// values given at the command line is correct.
func TestQuotaValueRx(t *testing.T) {
	tt := []struct {
		in    string
		match bool
	}{
		{"service/resource=123456", true},
		{"service/resource=123456:comment", true},
		{"service/resource=123.456Unit", true},
		{"service/resource=.456Unit:comment", true},

		{"serv1ce/resource=123", false},
		{"service/re3ource=123", false},
		{"service?resource=123", false},
		{"service/resource?123", false},
		{"service/resource=g123", false},
		{"service/resource=123g456", false},
		{"service/resource=123,456Unit", false},
	}

	for _, tc := range tt {
		if match := quotaValueRx.MatchString(tc.in); match != tc.match {
			if tc.match {
				t.Errorf("%q did not match the regular expression. Was expected to match.\n", tc.in)
			} else {
				t.Errorf("%q matched the regular expression. Was expected not to match.\n", tc.in)
			}
		}
	}
}

// TestParseCapacities tests if the service capacities given at the
// command line are being correctly parsed.
func TestParseCapacities(t *testing.T) {
	mockRawCapacities := RawQuotas{
		"shared/capacity=1.1MiB",
		"unshared/things=120:test comment.",
	}

	c, err := makeMockCluster("./fixtures/cluster-get-west.json")
	th.AssertNoErr(t, err)

	q, err := ParseRawQuotas(nil, c, &mockRawCapacities, true)
	th.AssertNoErr(t, err)

	actual := makeServiceCapacities(q)

	unitB := limes.UnitBytes
	unitNone := limes.UnitNone
	expected := []limes.ServiceCapacityRequest{
		{Type: "shared", Resources: []limes.ResourceCapacityRequest{
			{
				Name:     "capacity",
				Capacity: 1153433,
				Unit:     &unitB,
			},
		}},
		{Type: "unshared", Resources: []limes.ResourceCapacityRequest{
			{
				Name:     "things",
				Capacity: 120,
				Unit:     &unitNone,
				Comment:  "test comment.",
			},
		}},
	}
	th.AssertDeepEquals(t, expected, actual)
}

// TestParseQuotas tests if the service quotas given at the command line
// are being correctly parsed.
func TestParseQuotas(t *testing.T) {
	mockRawQuotas := RawQuotas{
		"shared/capacity=1.8KiB:this comment should be ignored by the parser.",
		"unshared/things=12",
	}

	assertParseQuotas := func(s baseUnitsSetter, rq *RawQuotas) {
		q, err := ParseRawQuotas(nil, s, rq, true)
		th.AssertNoErr(t, err)

		actual := makeServiceQuotas(q)

		expected := limes.QuotaRequest{
			"shared": limes.ServiceQuotaRequest{
				"capacity": limes.ValueWithUnit{1843, limes.UnitBytes},
			},
			"unshared": limes.ServiceQuotaRequest{
				"things": limes.ValueWithUnit{12, limes.UnitNone},
			},
		}
		th.AssertDeepEquals(t, expected, actual)
	}

	d, err := makeMockDomain("./fixtures/domain-get-germany.json")
	th.AssertNoErr(t, err)
	assertParseQuotas(d, &mockRawQuotas)

	p, err := makeMockProject("./fixtures/project-get-dresden.json")
	th.AssertNoErr(t, err)
	assertParseQuotas(p, &mockRawQuotas)
}
