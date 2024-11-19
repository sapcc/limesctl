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
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/utils/v2/client"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	"github.com/sapcc/go-bits/secrets"
	"github.com/sapcc/gophercloud-sapcc/v2/clients"
	"github.com/spf13/cobra"

	"github.com/sapcc/limesctl/v3/internal/util"
)

type VersionInfo struct {
	Version       string
	GitCommitHash string
	BuildDate     string
}

func Execute(ctx context.Context, v *VersionInfo) {
	if err := newRootCmd(v).ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Global flags.
var (
	debug bool

	osAuthURL           string
	osUsername          string
	osPassword          string
	osPwCmd             string
	osUserDomainID      string
	osUserDomainName    string
	osProjectID         string
	osProjectName       string
	osProjectDomainID   string
	osProjectDomainName string
	osCert              string
	osKey               string
)

func newRootCmd(v *VersionInfo) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "limesctl",
		Short: "Command-line client for Limes",
		Args:  cobra.NoArgs,
		Version: fmt.Sprintf("%s, Git commit %s, built at %s",
			v.Version, v.GitCommitHash, v.BuildDate),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Flags
	doNotSortFlags(cmd)
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug mode (will print API requests and responses)")
	cmd.PersistentFlags().StringVar(&osAuthURL, "os-auth-url", "", "authentication URL")
	cmd.PersistentFlags().StringVar(&osUsername, "os-username", "", "username")
	cmd.PersistentFlags().StringVar(&osPassword, "os-password", "", "user's Password")
	cmd.PersistentFlags().StringVar(&osPwCmd, "os-pw-cmd", "", "command from which to retrieve the user's password")
	cmd.PersistentFlags().StringVar(&osUserDomainID, "os-user-domain-id", "", "user's domain ID")
	cmd.PersistentFlags().StringVar(&osUserDomainName, "os-user-domain-name", "", "user's domain name")
	cmd.PersistentFlags().StringVar(&osProjectID, "os-project-id", "", "project ID to scope to")
	cmd.PersistentFlags().StringVar(&osProjectName, "os-project-name", "", "project name to scope to")
	cmd.PersistentFlags().StringVar(&osProjectDomainID, "os-project-domain-id", "", "domain ID containing project to scope to")
	cmd.PersistentFlags().StringVar(&osProjectDomainName, "os-project-domain-name", "", "domain name containing project to scope to")
	cmd.PersistentFlags().StringVar(&osCert, "os-cert", "", "client certificate")
	cmd.PersistentFlags().StringVar(&osKey, "os-key", "", "client certificate key")

	// Subcommands
	cmd.AddCommand(newClusterCmd())
	cmd.AddCommand(newDomainCmd())
	cmd.AddCommand(newProjectCmd())
	cmd.AddCommand(newOpsCmd())
	cmd.AddCommand(newLiquidCmd())

	return cmd
}

// Service clients that are used by different commands.
var (
	identityClient       *gophercloud.ServiceClient
	limesResourcesClient *gophercloud.ServiceClient
	limesRatesClient     *gophercloud.ServiceClient
	limesAdminClient     *gophercloud.ServiceClient
)

func authenticate(ctx context.Context) (*gophercloud.ProviderClient, error) {
	// Update OpenStack environment variables, if value(s) provided as flag.
	updateOpenStackEnvVars()

	if err := secrets.GetPasswordFromCommandIfRequested(); err != nil {
		return nil, err
	}
	ao, err := clientconfig.AuthOptions(nil)
	if err != nil {
		return nil, util.WrapError(err, "could not get auth variables")
	}

	provider, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return nil, util.WrapError(err, "cannot create an OpenStack client")
	}

	transport := &http.Transport{}
	if os.Getenv("OS_CERT") != "" && os.Getenv("OS_KEY") != "" {
		cert, err := tls.LoadX509KeyPair(os.Getenv("OS_CERT"), os.Getenv("OS_KEY"))
		if err != nil {
			return nil, util.WrapError(err, "failed to load x509 key pair")
		}
		transport.TLSClientConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
		provider.HTTPClient = http.Client{
			Transport: transport,
		}
	}

	if debug {
		provider.HTTPClient = http.Client{
			Transport: &client.RoundTripper{
				Rt:     transport,
				Logger: &client.DefaultLogger{},
			},
		}
	}

	err = openstack.Authenticate(ctx, provider, *ao)
	if err != nil {
		return nil, util.WrapError(err, "cannot connect to OpenStack")
	}

	identityClient, err = openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, util.WrapError(err, "could not initialize identity client")
	}

	return provider, nil
}

func authWithLimesResources(cmd *cobra.Command, _ []string) error {
	provider, err := authenticate(cmd.Context())
	if err != nil {
		return err
	}
	limesResourcesClient, err = clients.NewLimesV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return util.WrapError(err, "could not initialize Limes resources client")
	}
	return nil
}

func authWithLimesRates(cmd *cobra.Command, _ []string) error {
	provider, err := authenticate(cmd.Context())
	if err != nil {
		return err
	}
	limesRatesClient, err = clients.NewLimesRatesV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return util.WrapError(err, "could not initialize Limes rates client")
	}
	return nil
}

func authWithLimesAdmin(cmd *cobra.Command, _ []string) error {
	provider, err := authenticate(cmd.Context())
	if err != nil {
		return err
	}
	endpointOpts := gophercloud.EndpointOpts{}
	endpointOpts.ApplyDefaults("resources")
	endpoint, err := provider.EndpointLocator(endpointOpts)
	if err != nil {
		return util.WrapError(err, "could not initialize Limes admin client")
	}
	endpoint += "admin/"
	limesAdminClient = &gophercloud.ServiceClient{
		ProviderClient: provider,
		Endpoint:       endpoint,
	}
	return nil
}

func setenvIfVal(key, val string) {
	if val == "" {
		return
	}
	err := os.Setenv(key, val)
	if err != nil {
		cobra.CheckErr(err)
	}
}

func updateOpenStackEnvVars() {
	setenvIfVal("OS_AUTH_URL", osAuthURL)
	setenvIfVal("OS_USERNAME", osUsername)
	setenvIfVal("OS_PASSWORD", osPassword)
	setenvIfVal("OS_PW_CMD", osPwCmd)
	setenvIfVal("OS_USER_DOMAIN_ID", osUserDomainID)
	setenvIfVal("OS_USER_DOMAIN_NAME", osUserDomainName)
	setenvIfVal("OS_PROJECT_ID", osProjectID)
	setenvIfVal("OS_PROJECT_NAME", osProjectName)
	setenvIfVal("OS_PROJECT_DOMAIN_ID", osProjectDomainID)
	setenvIfVal("OS_PROJECT_DOMAIN_NAME", osProjectDomainName)
	setenvIfVal("OS_CERT", osCert)
	setenvIfVal("OS_KEY", osKey)
}
