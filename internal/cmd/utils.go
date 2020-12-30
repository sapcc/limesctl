// Copyright 2020 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/sapcc/limesctl/internal/core"
)

func writeJSON(d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return errors.Wrap(err, "could not marshal JSON")
	}

	if _, err = fmt.Fprintln(os.Stdout, string(b)); err != nil {
		return errors.Wrap(err, "could not write JSON data")
	}

	return nil
}

func writeReports(opts outputFormatFlags, reports ...core.LimesReportRenderer) error {
	d := core.RenderReports(opts.csvRecFmt, opts.Humanize, reports...)
	var err error
	if opts.Format == core.OutputFormatCSV {
		err = d.Write(os.Stdout)
	} else {
		d.WriteAsTable()
	}
	return err
}
