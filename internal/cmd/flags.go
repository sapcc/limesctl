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
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type RequestFilterFlags struct {
	Area     string `help:"Resource area."`
	Service  string `help:"Service type."`
	Resource string `help:"resource name."`
}

type OutputFormatFlags struct {
	Format        string `enum:"${outputFormats}" default:"table" help:"Output format (${enum})."`
	HumanReadable bool   `help:"Show quota and usage values in an user friendly unit."`
	Long          bool   `help:"Show detailed output."`
	Names         bool   `help:"Show output with names instead of UUIDs."`
}

type Globals struct {
	Debug   bool        `env:"LIMESCTL_DEBUG" help:"Enable debug mode."`
	Version VersionFlag `help:"Print version information and quit."`
	openStackFlags
}

// VersionFlag is a custom implementation of kong.VersionFlag.
// It is used to display the version info.
type VersionFlag struct {
	Version       string
	GitCommitHash string
	BuildDate     string
}

// Decode implements the kong.MapperValue interface.
func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }

// IsBool implements the kong.BoolMapper interface.
func (v VersionFlag) IsBool() bool { return true }

// BeforeApply writes the version info and terminates with a 0 exit status.
func (v VersionFlag) BeforeApply(app *kong.Kong, version VersionFlag) error {
	fmt.Printf("limesctl has version %s built from Git commit %s on %s\n",
		version.Version, version.GitCommitHash, version.BuildDate)
	app.Exit(0)
	return nil
}

type openStackFlags struct {
	OSUsername          string `help:"Username."`
	OSPassword          string `help:"User's Password."`
	OSUserDomainID      string `help:"User's domain ID."`
	OSUserDomainName    string `help:"User's domain name."`
	OSProjectID         string `help:"Project ID to scope to."`
	OSProjectName       string `help:"Project name to scope to."`
	OSProjectDomainID   string `help:"Domain ID containing project to scope to."`
	OSProjectDomainName string `help:"Domain name containing project to scope to."`
	OSAuthURL           string `help:"Authentication URL."`
}

func (o *openStackFlags) AfterApply() error {
	// Overwrite OpenStack environment variables
	if err := setEnvUnlessEmpty("OS_USERNAME", o.OSUsername); err != nil {
		return err
	}
	if err := setEnvUnlessEmpty("OS_PASSWORD", o.OSPassword); err != nil {
		return err
	}
	if err := setEnvUnlessEmpty("OS_USER_DOMAIN_ID", o.OSUserDomainID); err != nil {
		return err
	}
	if err := setEnvUnlessEmpty("OS_USER_DOMAIN_NAME", o.OSUserDomainName); err != nil {
		return err
	}
	if err := setEnvUnlessEmpty("OS_PROJECT_ID", o.OSProjectID); err != nil {
		return err
	}
	if err := setEnvUnlessEmpty("OS_PROJECT_NAME", o.OSProjectName); err != nil {
		return err
	}
	if err := setEnvUnlessEmpty("OS_PROJECT_DOMAIN_ID", o.OSProjectDomainID); err != nil {
		return err
	}
	if err := setEnvUnlessEmpty("OS_PROJECT_DOMAIN_NAME", o.OSProjectDomainName); err != nil {
		return err
	}
	if err := setEnvUnlessEmpty("OS_AUTH_URL", o.OSAuthURL); err != nil {
		return err
	}
	return nil
}

func setEnvUnlessEmpty(key, val string) error {
	if val != "" {
		if err := os.Setenv(key, val); err != nil {
			return err
		}
	}
	return nil
}
