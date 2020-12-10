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

package core

import "time"

func timestampToString(timestamp *int64) string {
	if timestamp == nil {
		return ""
	}
	return time.Unix(*timestamp, 0).UTC().Format(time.RFC3339)
}

func valFromPtr(ptr *uint64) uint64 {
	var v uint64
	if ptr != nil {
		v = *ptr
	}
	return v
}

func emptyStrIfZero(s string) string {
	if s == "0" {
		s = ""
	}
	return s
}
