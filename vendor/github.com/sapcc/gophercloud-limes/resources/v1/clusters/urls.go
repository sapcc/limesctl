package clusters

import "github.com/gophercloud/gophercloud"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("clusters")
}

func getURL(client *gophercloud.ServiceClient, clusterID string) string {
	return client.ServiceURL("clusters", clusterID)
}

func updateURL(client *gophercloud.ServiceClient, clusterID string) string {
	return client.ServiceURL("clusters", clusterID)
}
