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
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/utils/client"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/pkg/errors"
	"github.com/sapcc/gophercloud-sapcc/clients"
)

type CLI struct {
	globalFlags

	Cluster clusterCmd `cmd:"" help:"Do some action on cluster."`
	Domain  domainCmd  `cmd:"" help:"Do some action on domain(s)."`
	Project projectCmd `cmd:"" help:"Do some action on project(s)."`
}

// globalFlags holds app level (global) flags.
type globalFlags struct {
	Debug   bool        `env:"LIMESCTL_DEBUG" help:"Enable debug mode (will print API requests and responses)."`
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
	fmt.Printf("limesctl version %s built from Git commit %s on %s\n",
		version.Version, version.GitCommitHash, version.BuildDate)
	app.Exit(0)
	return nil
}

// openstackFlags holds the values for the required.
type openStackFlags struct {
	OSAuthURL           string `help:"Authentication URL."`
	OSUsername          string `help:"Username."`
	OSPassword          string `help:"User's Password."`
	OSUserDomainID      string `help:"User's domain ID."`
	OSUserDomainName    string `help:"User's domain name."`
	OSProjectID         string `help:"Project ID to scope to."`
	OSProjectName       string `help:"Project name to scope to."`
	OSProjectDomainID   string `help:"Domain ID containing project to scope to."`
	OSProjectDomainName string `help:"Domain name containing project to scope to."`
}

// ServiceClients holds the service clients for v3 identity service and Limes.
type ServiceClients struct {
	identity *gophercloud.ServiceClient
	limes    *gophercloud.ServiceClient
}

// Authenticate authenticates against OpenStack and returns the necessary
// service clients.
func (cli *CLI) Authenticate() (*ServiceClients, error) {
	// Update OpenStack environment variables, if value provided as flag.
	err := updateOpenStackEnvVars(&cli.openStackFlags)
	if err != nil {
		return nil, err
	}

	ao, err := clientconfig.AuthOptions(nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not get auth variables")
	}

	provider, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create an OpenStack client")
	}
	if cli.Debug {
		provider.HTTPClient = http.Client{
			Transport: &client.RoundTripper{
				Rt:     &http.Transport{},
				Logger: &client.DefaultLogger{},
			},
		}
	}

	err = openstack.Authenticate(provider, *ao)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to OpenStack")
	}

	identityClient, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize identity client")
	}

	limesClient, err := clients.NewLimesV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize Limes client")
	}

	return &ServiceClients{
		identity: identityClient,
		limes:    limesClient,
	}, nil
}

func setenvIfVal(key, val string) error {
	if val == "" {
		return nil
	}
	return os.Setenv(key, val)
}

func updateOpenStackEnvVars(v *openStackFlags) error {
	if err := setenvIfVal("OS_AUTH_URL", v.OSAuthURL); err != nil {
		return err
	}
	if err := setenvIfVal("OS_USERNAME", v.OSUsername); err != nil {
		return err
	}
	if err := setenvIfVal("OS_PASSWORD", v.OSPassword); err != nil {
		return err
	}
	if err := setenvIfVal("OS_USER_DOMAIN_ID", v.OSUserDomainID); err != nil {
		return err
	}
	if err := setenvIfVal("OS_USER_DOMAIN_NAME", v.OSUserDomainName); err != nil {
		return err
	}
	if err := setenvIfVal("OS_PROJECT_ID", v.OSProjectID); err != nil {
		return err
	}
	if err := setenvIfVal("OS_PROJECT_NAME", v.OSProjectName); err != nil {
		return err
	}
	if err := setenvIfVal("OS_PROJECT_DOMAIN_ID", v.OSProjectDomainID); err != nil {
		return err
	}
	if err := setenvIfVal("OS_PROJECT_DOMAIN_NAME", v.OSProjectDomainName); err != nil {
		return err
	}
	return nil
}
