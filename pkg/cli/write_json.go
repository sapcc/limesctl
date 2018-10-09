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
	"encoding/json"
	"fmt"

	"github.com/sapcc/go-bits/logg"
)

// writeJSON renders the result of a get/list/update in the JSON format
// and writes it to os.Stdout.
func (c *Cluster) writeJSON() {
	b, err := json.Marshal(c.Result.Body)
	if err != nil {
		logg.Fatal(err.Error())
	}

	fmt.Println(string(b))
}

// writeJSON renders the result of a get/list/update in the JSON format
// and writes it to os.Stdout.
func (d *Domain) writeJSON() {
	b, err := json.Marshal(d.Result.Body)
	if err != nil {
		logg.Fatal(err.Error())
	}

	fmt.Println(string(b))
}

// writeJSON renders the result of a get/list/update in the JSON format
// and writes it to os.Stdout.
func (p *Project) writeJSON() {
	b, err := json.Marshal(p.Result.Body)
	if err != nil {
		logg.Fatal(err.Error())
	}

	fmt.Println(string(b))
}
