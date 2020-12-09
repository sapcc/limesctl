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
