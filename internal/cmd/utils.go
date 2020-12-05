package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
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

func writeReports(opts outputFormatFlags, reports []core.LimesReportRenderer) error {
	d := core.RenderReports(reports, opts.csvRecFmt, opts.Humanize)
	var err error
	if opts.Format == core.OutputFormatCSV {
		err = writeCSV(os.Stdout, d)
	} else {
		writeTable(d)
	}
	return err
}

// writeCSV is used in unit tests therefore it has io.Writer as a parameter.
func writeCSV(w io.Writer, d core.CSVRecords) error {
	csvW := csv.NewWriter(w)
	csvW.Comma = rune(';') // Use semicolon as delimiter
	if err := csvW.WriteAll(d); err != nil {
		return errors.Wrap(err, "could not write CSV data")
	}
	return nil
}

func writeTable(d core.CSVRecords) {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader(d[0])
	t.AppendBulk(d[1:])
	t.Render()
}
