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
	"errors"

	"github.com/sapcc/limesctl/internal/core"
)

// resourceFilterFlags define parameters for Limes API requests that concern
// resources.
type resourceFilterFlags struct {
	Service  string `help:"Service type."`
	Area     string `help:"Resource area."`
	Resource string `help:"Resource name."`
}

// rateFilterFlags define parameters for Limes API requests that concern rates.
type rateFilterFlags struct {
	Service string `help:"Service type."`
}

type commonOutputFmtFlags struct {
	Format core.OutputFormat `short:"f" enum:"${outputFormats}" default:"table" help:"Output format (${enum})."`
	Names  bool              `help:"Show output with names instead of UUIDs. Not valid for 'json' output format."`

	// This is set by the corresponding validate().
	csvRecFmt core.CSVRecordFormat `kong:"-"`
}

// resourceOutputFmtFlags define how the app will print resource data.
type resourceOutputFmtFlags struct {
	commonOutputFmtFlags
	Long     bool `help:"Show detailed output. Not valid for 'json' output format."`
	Humanize bool `help:"Show quota and usage values in an user friendly unit. Not valid for 'json' output format."`
}

func (o *resourceOutputFmtFlags) validate() error {
	if o.Long && o.Names {
		return errors.New("'--long' and '--names' flags are mutually exclusive")
	}
	if o.Long {
		o.csvRecFmt = core.CSVRecordFormatLong
	}
	if o.Names {
		o.csvRecFmt = core.CSVRecordFormatNames
	}
	return nil
}

// rateOutputFmtFlags define how the app will print rate data.
type rateOutputFmtFlags struct {
	commonOutputFmtFlags
}

func (o *rateOutputFmtFlags) validate() error {
	if o.Names {
		o.csvRecFmt = core.CSVRecordFormatNames
	}
	return nil
}
