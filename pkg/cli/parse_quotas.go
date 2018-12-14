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
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/sapcc/limes"
)

// Set implements the kingpin.Value interface.
func (q *Quotas) Set(value string) error {
	value = strings.TrimSpace(value)
	// tmp contains the different components of a single parsed quota value. This makes it easier to refer
	// to individual components and pass them to the cli.Quotas map
	var tmp Resource

	// separate the quota value from its identifier
	idfVal := strings.SplitN(value, "=", 2)
	if len(idfVal) != 2 {
		return fmt.Errorf("expected a quota in the format: service/resource=value(unit), got '%s'", value)
	}

	// separate service and resource
	srvRes := strings.SplitN(idfVal[0], "/", 2)
	if len(srvRes) != 2 {
		return fmt.Errorf("expected service/resource, got '%s'", idfVal[0])
	}
	srv := srvRes[0]
	tmp.Name = srvRes[1]

	// separate quota value and comment (if one was given)
	valCom := strings.SplitN(idfVal[1], ":", 2)
	if len(valCom) > 1 {
		if valCom[1] != "" {
			tmp.Comment = strings.TrimSpace(valCom[1])
		}
	}

	// separate quota's value from its unit (if one was given)
	rx := regexp.MustCompile(`^([0-9]+)([A-Za-z]+)?$`)
	match := rx.MatchString(valCom[0])
	if !match {
		return fmt.Errorf("expected a quota value with optional unit in the format: 123Unit, got '%s'", valCom[0])
	}

	// rxMatchedList: []string{"entire regex matched string", "quota value", "unit (empty, if no unit given)"}
	rxMatchedList := rx.FindStringSubmatch(valCom[0])
	intVal, err := strconv.ParseInt(rxMatchedList[1], 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse quota value: '%s'", rxMatchedList[1])
	}
	tmp.Value = intVal

	tmp.Unit = limes.UnitNone
	if rxMatchedList[2] != "" {
		switch rxMatchedList[2] {
		case "B":
			tmp.Unit = limes.UnitBytes
		case "KiB":
			tmp.Unit = limes.UnitKibibytes
		case "MiB":
			tmp.Unit = limes.UnitMebibytes
		case "GiB":
			tmp.Unit = limes.UnitGibibytes
		case "TiB":
			tmp.Unit = limes.UnitTebibytes
		case "PiB":
			tmp.Unit = limes.UnitPebibytes
		case "EiB":
			tmp.Unit = limes.UnitExbibytes
		default:
			return fmt.Errorf("acceptable units: ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB'], got '%s'", rxMatchedList[2])
		}
	}

	(*q)[srv] = append((*q)[srv], tmp)

	return nil
}

// String implements the kingpin.Value interface.
func (q *Quotas) String() string {
	return ""
}

// IsCumulative allows consumption of remaining command line arguments.
func (q *Quotas) IsCumulative() bool {
	return true
}

// ParseQuotas parses a command line argument to a quota value and assigns it to the
// aggregate cli.Quotas map.
func ParseQuotas(s kingpin.Settings) (target *Quotas) {
	target = &Quotas{}
	s.SetValue((*Quotas)(target))
	return
}

// makeServiceCapacities is a helper function that converts a Quotas type to
// a slice of limes.ServiceCapacityRequest for use with cluster set operations.
func makeServiceCapacities(q *Quotas) []limes.ServiceCapacityRequest {
	//serialize service types with ordered keys for consistent test results
	types := make([]string, 0, len(*q))
	for typeStr := range *q {
		types = append(types, typeStr)
	}
	sort.Strings(types)

	sc := make([]limes.ServiceCapacityRequest, 0, len(types))
	for _, srv := range types {
		resList := (*q)[srv]
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
func makeServiceQuotas(q *Quotas) limes.QuotaRequest {
	sq := make(limes.QuotaRequest)

	for srv, resList := range *q {
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
