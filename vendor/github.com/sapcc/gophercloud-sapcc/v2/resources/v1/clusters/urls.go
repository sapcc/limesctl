// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package clusters

import "github.com/gophercloud/gophercloud/v2"

func getURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("clusters", "current")
}
