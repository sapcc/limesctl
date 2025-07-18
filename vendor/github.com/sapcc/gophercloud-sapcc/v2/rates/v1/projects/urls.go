// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package projects

import "github.com/gophercloud/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient, domainID string) string {
	return client.ServiceURL("domains", domainID, "projects")
}

func getURL(client *gophercloud.ServiceClient, domainID, projectID string) string {
	return client.ServiceURL("domains", domainID, "projects", projectID)
}
