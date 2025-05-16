// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

// token contains domain and project information about an authorization token.
type token struct {
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

// currentToken returns the current auth token that was used to authenticate
// against an OpenStack cloud.
func currentToken(identityClient *gophercloud.ServiceClient) (*token, error) {
	var t token
	err := identityClient.GetAuthResult().(tokens.CreateResult).ExtractInto(&t)
	if err != nil {
		return nil, fmt.Errorf("could not get current token: %w", err)
	}
	return &t, nil
}
