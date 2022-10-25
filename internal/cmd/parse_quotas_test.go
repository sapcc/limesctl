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

package cmd

import (
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/sapcc/go-api-declarations/limes"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
)

// TestSplitQuotaRe tests the `splitQuotaRe` regular expression.
func TestSplitQuotaRe(t *testing.T) {
	tt := []struct {
		in          string
		shouldMatch bool
	}{
		{"service/resource=123.456Unit", true},
		{"service/resource-=123.456Unit", true},
		{"service/resource+=123.456Unit", true},
		{"service/resource*=123.456Unit", true},
		{"service/resource/=123.456Unit", true},

		{"service/resource:123.456Unit", false},
		{"service/resource+123.456Unit", false},
		{"service/resource-123.456Unit", false},
		{"service/resource*123.456Unit", false},
		{"service/resource/123.456Unit", false},
	}

	for _, tc := range tt {
		if match := splitQuotaRe.MatchString(tc.in); match != tc.shouldMatch {
			if tc.shouldMatch {
				t.Errorf("%q did not match the regular expression. Was expected to match.\n", tc.in)
			} else {
				t.Errorf("%q matched the regular expression. Was expected not to match.\n", tc.in)
			}
		}
	}
}

// TestQuotaValueRe tests the `quotaValueRe` regular expression.
func TestQuotaValueRe(t *testing.T) {
	tt := []struct {
		in          string
		shouldMatch bool
	}{
		{"123456", true},
		{"123456Unit", true},
		{"123.456", true},
		{"123.456Unit", true},
		{"0.456", true},
		{"0.456Unit", true},
		{".456", true},
		{".456Unit", true},

		{"123456Un1t", false},
		{"123456Un-t", false},
		{"123a456Unit", false},
		{"123456-Unit", false},
		{"+123456Unit", false},
	}

	for _, tc := range tt {
		if match := quotaValueRe.MatchString(tc.in); match != tc.shouldMatch {
			if tc.shouldMatch {
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
	expected := limesresources.QuotaRequest{
		"shared": limesresources.ServiceQuotaRequest{
			"capacity":  limesresources.ResourceQuotaRequest{Value: 1597, Unit: limes.UnitMebibytes},
			"capacity2": limesresources.ResourceQuotaRequest{Value: 30, Unit: limes.UnitGibibytes},
		},
		"unshared": limesresources.ServiceQuotaRequest{
			"things":  limesresources.ResourceQuotaRequest{Value: 12, Unit: limes.UnitNone},
			"things2": limesresources.ResourceQuotaRequest{Value: 20, Unit: limes.UnitNone},
			"things3": limesresources.ResourceQuotaRequest{Value: 4, Unit: limes.UnitNone},
		},
	}

	resQuotas := resourceQuotas{
		"shared": map[string]limes.ValueWithUnit{
			"capacity":  {Value: 1024, Unit: limes.UnitMebibytes},
			"capacity2": {Value: 3, Unit: limes.UnitGibibytes},
		},
		"unshared": map[string]limes.ValueWithUnit{
			"things":  {Value: 10, Unit: limes.UnitNone},
			"things2": {Value: 10, Unit: limes.UnitNone},
			"things3": {Value: 44, Unit: limes.UnitNone},
		},
	}

	successMock := []string{
		"shared/capacity+=.56GiB",
		"shared/capacity2*=10.2567GiB",
		"unshared/things=12",
		"unshared/things2+=10",
		"unshared/things3/=10",
	}
	actual, err := parseToQuotaRequest(resQuotas, successMock)
	if err != nil {
		t.Error(err)
	}
	th.AssertDeepEquals(t, expected, actual)

	failureMocks := [][]string{
		{"unshared/things=12.12"}, // counted resources need to be an integer
		{"unshared/things2-=12"},  // will result in a quota value < 0
		{"unshared/things3/=45"},  // will result in a quota value < 0
	}
	for _, mock := range failureMocks {
		_, err := parseToQuotaRequest(resQuotas, mock)
		if err == nil {
			t.Error("expected an error, got nil.")
		}
	}
}
