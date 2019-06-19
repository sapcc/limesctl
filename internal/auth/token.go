package auth

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

// Token contains domain and project information about an authorisation token.
type Token struct {
	Domain struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"domain"`
	Project struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Domain struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"domain"`
	} `json:"project"`
}

// CurrentToken returns the current auth token that was used to authenticate
// against an OpenStack cloud.
func CurrentToken(identityClient *gophercloud.ServiceClient) (*Token, error) {
	var t Token

	err := identityClient.GetAuthResult().(tokens.CreateResult).ExtractInto(&t)
	if err != nil {
		return nil, fmt.Errorf("could not get current token: %v", err)
	}

	return &t, nil
}
