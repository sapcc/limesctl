// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package domains

import "github.com/gophercloud/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("domains")
}

func getURL(client *gophercloud.ServiceClient, domainID string) string {
	return client.ServiceURL("domains", domainID)
}
