// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sapcc/limesctl/v3/internal/core"
	"github.com/sapcc/limesctl/v3/internal/util"
)

func writeJSON(d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return util.WrapError(err, "could not marshal JSON")
	}

	if _, err = fmt.Fprintln(os.Stdout, string(b)); err != nil {
		return util.WrapError(err, "could not write JSON data")
	}

	return nil
}

func writeReports(opts *core.OutputOpts, reports ...core.LimesReportRenderer) error {
	d := core.RenderReports(opts, reports...)
	var err error
	if opts.Fmt == core.OutputFormatCSV {
		err = d.Write(os.Stdout)
	} else {
		d.WriteAsTable()
	}
	return err
}
