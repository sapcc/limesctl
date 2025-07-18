// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

// Package projects provides interaction with Limes at the project hierarchical
// level.
//
// Here is an example on how you would list all the projects in the current
// domain:
//
//	import (
//	  "fmt"
//	  "log"
//
//	  "github.com/gophercloud/gophercloud/v2"
//	  "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
//	  "github.com/gophercloud/utils/v2/openstack/clientconfig"
//
//	  "github.com/sapcc/gophercloud-sapcc/v2/clients"
//	  "github.com/sapcc/gophercloud-sapcc/v2/resources/v1/projects"
//	)
//
//	func main() {
//	  provider, err := clientconfig.AuthenticatedClient(nil)
//	  if err != nil {
//	    log.Fatalf("could not initialize openstack client: %v", err)
//	  }
//
//	  limesClient, err := clients.NewLimesV1(provider, gophercloud.EndpointOpts{})
//	  if err != nil {
//	    log.Fatalf("could not initialize Limes client: %v", err)
//	  }
//
//	  project, err := provider.GetAuthResult().(tokens.CreateResult).ExtractProject()
//	  if err != nil {
//	    log.Fatalf("could not get project from token: %v", err)
//	  }
//
//	  result := projects.List(limesClient, project.Domain.ID, projects.ListOpts{Detail: true})
//	  if result.Err != nil {
//	    log.Fatalf("could not get projects: %v", result.Err)
//	  }
//
//	  projectList, err := result.ExtractProjects()
//	  if err != nil {
//	    log.Fatalf("could not get projects: %v", err)
//	  }
//	  for _, project := range projectList {
//	    fmt.Printf("%+v\n", project.Services)
//	  }
//	}
package projects
