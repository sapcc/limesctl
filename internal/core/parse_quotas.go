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
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/gophercloud-limes/resources/v1/clusters"
	"github.com/sapcc/gophercloud-limes/resources/v1/domains"
	"github.com/sapcc/gophercloud-limes/resources/v1/projects"
	"github.com/sapcc/limes"
	"github.com/sapcc/limesctl/internal/errors"
)

// quotaUnits is map of limes.Unit to x, such that base-2 exponential of x is
// the number of bytes for some specific unit.
var quotaUnits = map[limes.Unit]float64{
	limes.UnitBytes:     0,
	limes.UnitKibibytes: 10,
	limes.UnitMebibytes: 20,
	limes.UnitGibibytes: 30,
	limes.UnitTebibytes: 40,
	limes.UnitPebibytes: 50,
	limes.UnitExbibytes: 60,
}

// RawQuotas contains the quota values provided at the command line.
type RawQuotas []string

// Resource contains quota information about a single resource.
type Resource struct {
	Name    string
	Value   int64
	Unit    limes.Unit
	Comment string
}

// Quotas is a map of service name to a list of resources. It contains the
// aggregate quota values used by the set methods to update a single
// cluster/domain/project.
type Quotas map[string][]Resource

// resourceUnits type contains the respective units for different resources.
type resourceUnits map[string]map[string]limes.Unit

// baseUnitsSetter is the interface type that is implemented by different
// hierarchies.
type baseUnitsSetter interface {
	setBaseUnits(*gophercloud.ServiceClient, *resourceUnits, bool)
}

// Reference:
//   matchList == [<full-match>, <service>, <resource>, <value>, <unit>?, <comment>?]
var quotaValueRx = regexp.MustCompile(`^([^:/=]+)/([^:/=]+)=(\d*\.?\d+)([a-zA-Z]+)?(?::(.*))?$`)

// ParseRawQuotas parses the raw quota values given at the command line to a
// Quotas map.
func ParseRawQuotas(limesV1 *gophercloud.ServiceClient, s baseUnitsSetter, rq RawQuotas, isTest bool) (Quotas, error) {
	q := make(Quotas)

	type userInput struct {
		service  string
		resource string
		value    string
		unit     limes.Unit
		comment  string
	}

	userInputs := make([]userInput, 0, len(rq))
	resUnits := make(resourceUnits)

	// validate raw quota values
	for _, rqInList := range rq {
		matchList := quotaValueRx.FindStringSubmatch(rqInList)
		if matchList == nil {
			return nil, fmt.Errorf("expected a quota with optional unit and comment in the format: service/resource=value(unit):comment, got %q", rqInList)
		}

		input := userInput{
			service:  matchList[1],
			resource: matchList[2],
			value:    matchList[3],
			comment:  matchList[5],
		}

		givenUnit := matchList[4]
		if givenUnit == "" {
			input.unit = limes.UnitNone
			if strings.Contains(input.value, ".") {
				return nil, fmt.Errorf("counted values must be an integer, got '%s'", input.value)
			}
		} else {
			input.unit = limes.Unit(givenUnit)
			_, unitIsValid := quotaUnits[input.unit]
			if !unitIsValid {
				acceptableUnits := make([]limes.Unit, 0, len(quotaUnits))
				for k := range quotaUnits {
					acceptableUnits = append(acceptableUnits, k)
				}
				sort.Slice(acceptableUnits, func(i, j int) bool {
					return quotaUnits[acceptableUnits[i]] < quotaUnits[acceptableUnits[j]]
				})
				return nil, fmt.Errorf("expected a unit from: %v, got %q", acceptableUnits, givenUnit)
			}
		}

		_, exists := resUnits[input.service]
		if !exists {
			resUnits[input.service] = make(map[string]limes.Unit)
		}
		resUnits[input.service][input.resource] = input.unit

		userInputs = append(userInputs, input)
	}

	s.setBaseUnits(limesV1, &resUnits, isTest)

	for _, input := range userInputs {
		var intQuotaVal int64

		if strings.Contains(input.value, ".") {
			diffInExp := quotaUnits[input.unit] - quotaUnits[resUnits[input.service][input.resource]]
			// fractional values are only possible when the given unit is
			// greater than the base unit
			if diffInExp < 0 {
				return nil, fmt.Errorf("minimum accepted unit for '%s/%s' is '%v', got '%v'",
					input.service, input.resource, resUnits[input.service][input.resource], input.unit)
			}

			// convert input.value to its base unit
			inputValue, err := strconv.ParseFloat(input.value, 64)
			if err != nil {
				return nil, err
			}
			inputValue = math.Floor(inputValue * math.Exp2(diffInExp))

			intQuotaVal = int64(inputValue)
			input.unit = resUnits[input.service][input.resource]
		} else {
			tmp, err := strconv.ParseInt(input.value, 10, 64)
			if err != nil {
				return nil, err
			}

			intQuotaVal = tmp
		}

		q[input.service] = append(q[input.service], Resource{
			Name:    input.resource,
			Value:   intQuotaVal,
			Unit:    input.unit,
			Comment: input.comment,
		})
	}

	return q, nil
}

// setBaseUnits (re)sets the resourceUnits map to the respective base units for
// resources.
func (c *Cluster) setBaseUnits(limesV1 *gophercloud.ServiceClient, ru *resourceUnits, isTest bool) {
	// no need to make a GET request if `isTest = true` because in that case
	// the c.Result would already be populated with the necessary mock test
	// data
	if !isTest {
		c.Result = clusters.Get(limesV1, c.ID, clusters.GetOpts{})
	}

	cluster, err := c.Result.Extract()
	errors.Handle(err, "cluster does not exist")

	for srv, resMap := range *ru {
		for res := range resMap {
			srvRes, exists := cluster.Services[srv].Resources[res]
			if exists {
				(*ru)[srv][res] = srvRes.ResourceInfo.Unit
			}
		}
	}
}

// setBaseUnits (re)sets the resourceUnits map to the respective base units for
// resources.
func (d *Domain) setBaseUnits(limesV1 *gophercloud.ServiceClient, ru *resourceUnits, isTest bool) {
	// no need to make a GET request if `isTest = true` because in that case
	// the d.Result would already be populated with the necessary mock test
	// data
	if !isTest {
		d.Result = domains.Get(limesV1, d.ID, domains.GetOpts{
			Cluster: d.Filter.Cluster})
	}

	domain, err := d.Result.Extract()
	errors.Handle(err, "domain does not exist")

	for srv, resMap := range *ru {
		for res := range resMap {
			srvRes, exists := domain.Services[srv].Resources[res]
			if exists {
				(*ru)[srv][res] = srvRes.ResourceInfo.Unit
			}
		}
	}
}

// setBaseUnits (re)sets the resourceUnits map to the respective base units for
// resources.
func (p *Project) setBaseUnits(limesV1 *gophercloud.ServiceClient, ru *resourceUnits, isTest bool) {
	// no need to make a GET request if `isTest = true` because in that case
	// the p.Result would already be populated with the necessary mock test
	// data
	if !isTest {
		p.Result = projects.Get(limesV1, p.DomainID, p.ID, projects.GetOpts{
			Cluster: p.Filter.Cluster})
	}

	project, err := p.Result.Extract()
	errors.Handle(err, "project does not exist")

	for srv, resMap := range *ru {
		for res := range resMap {
			srvRes, exists := project.Services[srv].Resources[res]
			if exists {
				(*ru)[srv][res] = srvRes.ResourceInfo.Unit
			}
		}
	}
}

// makeServiceCapacities is a helper function that converts a Quotas type to
// a slice of limes.ServiceCapacityRequest for use with cluster set operations.
func makeServiceCapacities(q Quotas) []limes.ServiceCapacityRequest {
	//serialize service types with ordered keys for consistent test results
	types := make([]string, 0, len(q))
	for typeStr := range q {
		types = append(types, typeStr)
	}
	sort.Strings(types)

	sc := make([]limes.ServiceCapacityRequest, 0, len(types))
	for _, srv := range types {
		resList := q[srv]
		rc := make([]limes.ResourceCapacityRequest, 0, len(resList))
		for _, r := range resList {
			// take a copy of the loop variable (it will be updated by the loop, so if
			// we didn't take a copy manually, the 'r' inside the list would
			// contain only identical pointers)
			r := r
			rc = append(rc, limes.ResourceCapacityRequest{
				Name:     r.Name,
				Capacity: r.Value,
				Unit:     &r.Unit,
				Comment:  r.Comment,
			})
		}

		sc = append(sc, limes.ServiceCapacityRequest{
			Type:      srv,
			Resources: rc,
		})
	}

	return sc
}

// makeServiceQuotas is a helper function that converts a Quotas type to
// limes.QuotaRequest for use with domain/project set operations.
func makeServiceQuotas(q Quotas) limes.QuotaRequest {
	sq := make(limes.QuotaRequest)

	for srv, resList := range q {
		sq[srv] = make(limes.ServiceQuotaRequest)

		for _, r := range resList {
			sq[srv][r.Name] = limes.ValueWithUnit{
				Value: uint64(r.Value),
				Unit:  r.Unit,
			}
		}
	}

	return sq
}
