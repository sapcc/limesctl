/*******************************************************************************
*
* Copyright 2018-2020 SAP SE
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

package cmd

import (
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/sapcc/limes"
)

// TestQuotaRx tests if the regular expression used to parse the quota values
// given at the command-line is correct.
func TestQuotaRx(t *testing.T) {
	tt := []struct {
		in    string
		match bool
	}{
		{"ser1vice/reso-urce=123456", true},
		{"service/resource=123.456Unit", true},

		{"serv/ice/resource=123", false},
		{"service/resou=rce=123", false},
		{"service/reso:urce=123", false},
		{"service?resource=123", false},
		{"service/resource?123", false},
		{"service/resource=g123", false},
		{"service/resource=123g456", false},
		{"service/resource=123,456Unit", false},
	}

	for _, tc := range tt {
		if match := quotaRx.MatchString(tc.in); match != tc.match {
			if tc.match {
				t.Errorf("%q did not match the regular expression. Was expected to match.\n", tc.in)
			} else {
				t.Errorf("%q matched the regular expression. Was expected not to match.\n", tc.in)
			}
		}
	}
}

// TestParseToQuotaRequest tests if the resource quotas given at the
// command-line are being correctly parsed.
func TestParseToQuotaRequest(t *testing.T) {
	expected := limes.QuotaRequest{
		"shared": limes.ServiceQuotaRequest{
			Resources: limes.ResourceQuotaRequest{
				"capacity":    limes.ValueWithUnit{Value: 1675037245, Unit: limes.UnitBytes},
				"capacityTwo": limes.ValueWithUnit{Value: 10, Unit: limes.UnitGibibytes},
			},
		},
		"unshared": limes.ServiceQuotaRequest{
			Resources: limes.ResourceQuotaRequest{
				"things": limes.ValueWithUnit{Value: 12, Unit: limes.UnitNone},
			},
		},
	}

	defaultResUnits := resourceUnits{
		"shared": map[string]limes.Unit{
			"capacity":    limes.UnitMebibytes,
			"capacityTwo": limes.UnitGibibytes,
		},
		"unshared": map[string]limes.Unit{
			"things": limes.UnitNone,
		},
	}
	mock := []string{
		"shared/capacity=1.56GiB",
		"shared/capacityTwo=6.6GiB",
		"unshared/things=12",
	}
	actual, err := parseToQuotaRequest(defaultResUnits, mock)
	if err != nil {
		t.Error(err)
	}

	th.AssertDeepEquals(t, expected, actual)
}
