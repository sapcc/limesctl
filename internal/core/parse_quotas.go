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

package core

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/limes"
)

// quotaRx is used to extract the quota information from user input. See
// parseToQuotaRequest() for more info.
// Reference:
//   matchList == [<full-match>, <service>, <resource>, <value>, <unit>?]
var quotaRx = regexp.MustCompile(`^([^:/=]+)/([^:/=]+)=(\d*\.?\d+)([a-zA-Z]+)?$`)

// parseToQuotaRequest parses a slice of user input quota values, converts the
// values to the resource's default unit (if needed), and returns a
// limes.QuotaRequest for use with PUT requests on domains and projects.
//
// The input values are expected to be in the format:
//   service/resource=123(Unit)
// where unit is optional.
func parseToQuotaRequest(resUnits ResourceUnits, in []string) (limes.QuotaRequest, error) {
	out := make(limes.QuotaRequest)
	for _, inputStr := range in {
		matchList := quotaRx.FindStringSubmatch(inputStr)
		if matchList == nil {
			return nil, fmt.Errorf("expected a quota/capacity with optional unit and comment in the format: service/resource=value(unit):comment, got %q", inputStr)
		}

		// Validate input service/resource.
		service := matchList[1]
		resource := matchList[2]
		defaultUnit, ok := resUnits[service][resource]
		if !ok {
			return nil, fmt.Errorf("invalid resource: %s/%s does not exist in Limes", service, resource)
		}

		valStr := matchList[3]
		isFloatVal := strings.Contains(valStr, ".")
		unitStr := matchList[4]

		// Validate input unit.
		var unit limes.Unit
		if unitStr == "" {
			if isFloatVal {
				return nil, fmt.Errorf("counted (i.e. resource without any unit) values must be an integer, got %q", valStr)
			}
			unit = limes.UnitNone
		} else {
			for _, v := range LimesUnits {
				if unitStr == string(v) {
					unit = v
					break
				}
			}
			if unit == "" {
				return nil, fmt.Errorf("expected a unit from %q, got %q", LimesUnits, unitStr)
			}
		}

		// Validate and convert (if needed) input value.
		var vWithUnit limes.ValueWithUnit
		if isFloatVal {
			logg.Info("Limes only accepts integer values, will attempt to convert %s %s to a suitable unit for %s/%s)",
				valStr, unit, service, resource)
			var err error
			vWithUnit, err = convertTo(valStr, unit, defaultUnit)
			if err != nil {
				if _, ok := err.(limes.IncompatibleUnitsError); ok {
					err = fmt.Errorf("can not convert %s %s, minimum accepted unit for %s/%s is %s",
						valStr, unit, service, resource, defaultUnit)
				}
				return nil, err
			}
		} else {
			v, err := strconv.ParseUint(valStr, 10, 64)
			if err != nil {
				return nil, err
			}
			vWithUnit = limes.ValueWithUnit{Value: v, Unit: unit}
		}

		if _, ok := out[service]; !ok {
			out[service] = limes.ServiceQuotaRequest{
				Resources: make(limes.ResourceQuotaRequest),
			}
		}
		out[service].Resources[resource] = vWithUnit
	}

	return out, nil
}

// convertTo converts the given value from source to target unit and returns
// the truncated result in limes.ValueWithUnit.
// i.e: 22.65 TiB -> 23193 GiB (instead of 23193.6 GiB)
func convertTo(valStr string, source, target limes.Unit) (limes.ValueWithUnit, error) {
	v, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return limes.ValueWithUnit{}, err
	}

	var newV float64
	if source == target {
		newV = math.Floor(v)
	} else {
		_, sourceMultiple := source.Base()
		_, targetMultiple := target.Base()
		if sourceMultiple > targetMultiple {
			return limes.ValueWithUnit{}, limes.IncompatibleUnitsError{Source: source, Target: target}
		}
		vInBase := math.Floor(v * float64(sourceMultiple))
		newV = math.Floor(vInBase / float64(targetMultiple))
	}

	logg.Info("%s %s -> %.0f %s", valStr, source, newV, target)
	return limes.ValueWithUnit{
		Value: uint64(newV),
		Unit:  target,
	}, nil
}
