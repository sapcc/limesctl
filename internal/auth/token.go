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

package auth

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
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
		return nil, fmt.Errorf("could not get current token: %v", err)
	}
	return &t, nil
}
