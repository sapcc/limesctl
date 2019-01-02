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
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/sapcc/limes"
)

var mockRawCapacities = []string{
	"compute/cores=10",
	"compute/ram=20MiB",
	"object-store/capacity=30B:I got 99 problems, but a cluster ain't one.",
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

var mockRawQuotas = []string{
	"compute/cores=10",
	"compute/ram=20MiB",
	"object-store/capacity=30B:this comment should be ignored by the parser.",
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
