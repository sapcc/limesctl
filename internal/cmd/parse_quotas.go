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
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sapcc/go-api-declarations/limes"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"

	"github.com/sapcc/limesctl/v3/internal/core"
)

// resourceQuotas is a map of service name to resource name to the resource's
// quota value and unit.
type resourceQuotas map[limes.ServiceType]map[limesresources.ResourceName]limes.ValueWithUnit

// splitQuotaRe is used to split the user input around the equality sign.
// Reference:
//
//	The following user input: service/resource+=123.456Unit
//	will result in:
//	  matchList == [<full-match>, "service/resource", "+=", "123.456Unit"]
var splitQuotaRe = regexp.MustCompile(`^(\S*?)([-+*/]?=)(\S*)$`)

// quotaValueRe is used to extract quota value and unit.
// Reference:
//
//	The following user input: 123.456Unit
//	will result in:
//	  matchList == [<full-match>, "123.456", "Unit"]
var quotaValueRe = regexp.MustCompile(`^(\d*\.?\d+)([a-zA-Z]*)$`)

// parseToQuotaRequest parses a slice of user input quota values, converts the
// values to the resource's default unit (if needed), and returns a
// limes.QuotaRequest for use with PUT requests on domains and projects.
//
// The input values are expected to be in the format:
//
//	service/resource=123(Unit)
//
// where unit is optional.
func parseToQuotaRequest(resValues resourceQuotas, in []string) (limesresources.QuotaRequest, error) {
	out := make(limesresources.QuotaRequest)
	for _, inputStr := range in {
		matchList := splitQuotaRe.FindStringSubmatch(inputStr)
		if matchList == nil {
			return out, fmt.Errorf("expected a quota with optional unit in the format: service/resource=value(unit), got %q", inputStr)
		}

		// Validate input service/resource.
		serviceResource := strings.Split(matchList[1], "/")
		if len(serviceResource) > 2 {
			return out, fmt.Errorf("expected a quota with service resource in the format: service/resource, got %q", matchList[1])
		}
		service := limes.ServiceType(serviceResource[0])
		resource := limesresources.ResourceName(serviceResource[1])
		currentValWithUnit, ok := resValues[service][resource]
		if !ok {
			return nil, fmt.Errorf("invalid resource: %s/%s does not exist in Limes", service, resource)
		}

		// Check if the same resource was given multiple times.
		if srv, ok := out[service]; ok {
			if _, ok := srv[resource]; ok {
				return nil, fmt.Errorf(
					"%s/%s was given multiple times, limesctl only supports a single change for a specific resource for a given request",
					service, resource)
			}
		}

		// Validate input value.
		valueUnitML := quotaValueRe.FindStringSubmatch(matchList[3])
		if valueUnitML == nil {
			return out, fmt.Errorf("expected a quota value with optional unit in the format: value(unit), got %q", matchList[3])
		}
		valStr := valueUnitML[1]
		isFloatVal := strings.Contains(valStr, ".")

		// Validate input unit.
		var unit limes.Unit
		if unitStr := valueUnitML[2]; unitStr == "" {
			if isFloatVal {
				return nil, fmt.Errorf("counted (i.e. resource without any unit) values must be an integer, got %q", valStr)
			}
			unit = limes.UnitNone
		} else {
			for _, v := range core.LimesUnits {
				if unitStr == string(v) {
					unit = v
					break
				}
			}
			if unit == "" {
				return nil, fmt.Errorf("expected a unit from %q, got %q", core.LimesUnits, unitStr)
			}
		}

		// Validate and convert (if needed) input value.
		operation := matchList[2]
		var newValWithUnit limes.ValueWithUnit
		if isFloatVal || operation != "=" {
			if isFloatVal {
				fmt.Fprintf(os.Stderr, "warning: Limes only accepts integer values, will attempt to convert %s %s to a suitable unit for %s/%s",
					valStr, unit, service, resource)
			}
			var err error
			newValWithUnit, err = convertTo(valStr, unit, currentValWithUnit.Unit)
			if err != nil {
				return nil, err
			}
			if isFloatVal {
				fmt.Fprintf(os.Stderr, "warning: converted %s %s -> %s", valStr, unit, newValWithUnit.String())
			}
		} else {
			v, err := strconv.ParseUint(valStr, 10, 64)
			if err != nil {
				return nil, err
			}
			newValWithUnit = limes.ValueWithUnit{Value: v, Unit: unit}
		}

		switch operation {
		case "+=":
			newValWithUnit.Value += currentValWithUnit.Value
		case "-=":
			if newValWithUnit.Value > currentValWithUnit.Value {
				return nil, fmt.Errorf("invalid quota value: subtraction of %s %s for %s/%s will result in a value < 0",
					valStr, unit, service, resource)
			}
			newValWithUnit.Value = currentValWithUnit.Value - newValWithUnit.Value
		case "*=":
			newValWithUnit.Value *= currentValWithUnit.Value
		case "/=":
			if newValWithUnit.Value > currentValWithUnit.Value {
				return nil, fmt.Errorf("invalid quota value: division by %s %s for %s/%s will result in a value < 1",
					valStr, unit, service, resource)
			}
			newValWithUnit.Value = currentValWithUnit.Value / newValWithUnit.Value
		}

		if _, ok := out[service]; !ok {
			out[service] = make(limesresources.ServiceQuotaRequest)
		}
		out[service][resource] = limesresources.ResourceQuotaRequest(newValWithUnit)
	}

	return out, nil
}

// smallerThanMinimumUnitError is returned by convertTo() when a value has a
// source unit that is smaller than the default unit for that particular
// resource.
type smallerThanDefaultUnitError struct {
	Value           string
	Source          limes.Unit
	ResourceDefault limes.Unit
}

// Error implements the error interface.
func (e smallerThanDefaultUnitError) Error() string {
	return fmt.Sprintf("cannot convert %s %s, minimum accepted unit for this resource is %s", e.Value, e.Source, e.ResourceDefault)
}

// convertTo converts the given value from source to target unit and returns
// the truncated result in limes.ValueWithUnit.
// i.e: 22.65 TiB -> 23193 GiB (instead of 23193.6 GiB).
func convertTo(valStr string, source, target limes.Unit) (limes.ValueWithUnit, error) {
	v, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return limes.ValueWithUnit{}, err
	}

	var newV float64
	if source == target {
		newV = math.Floor(v)
	} else {
		sourceBase, sourceMultiple := source.Base()
		targetBase, targetMultiple := target.Base()
		if sourceBase != targetBase {
			return limes.ValueWithUnit{}, limes.IncompatibleUnitsError{Source: source, Target: target}
		}
		if sourceMultiple < targetMultiple {
			return limes.ValueWithUnit{}, smallerThanDefaultUnitError{Value: valStr, Source: source, ResourceDefault: target}
		}
		vInBase := math.Floor(v * float64(sourceMultiple))
		newV = math.Floor(vInBase / float64(targetMultiple))
	}

	return limes.ValueWithUnit{
		Value: uint64(newV),
		Unit:  target,
	}, nil
}
