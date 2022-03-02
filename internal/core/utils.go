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

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	th "github.com/gophercloud/gophercloud/testhelper"
)

func timestampToString(timestamp *int64) string {
	if timestamp == nil {
		return ""
	}
	return time.Unix(*timestamp, 0).UTC().Format(time.RFC3339)
}

func zeroIfNil(ptr *uint64) uint64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func emptyStrIfNil(ptr *uint64, valToStr valToStrFunc) string {
	if ptr == nil {
		return ""
	}
	return valToStr(*ptr)
}

///////////////////////////////////////////////////////////////////////////////
// Helper functions for unit tests.

func fixturePath(filename string) string {
	return filepath.Join("fixtures", filename)
}

func fixtureBytes(filename string) ([]byte, error) {
	return os.ReadFile(filepath.Join("fixtures", filename))
}

// assertEquals is like testhelper.AssertEquals() but also writes actual
// content to file to make it easy to copy the computed result over to the
// fixture path when a new test is added or an existing one is modified.
func assertEquals(t *testing.T, fixtureFilename string, actual []byte) {
	t.Helper()
	fixturePathAbs, err := filepath.Abs(fixturePath(fixtureFilename))
	th.AssertNoErr(t, err)
	actualPathAbs := fixturePathAbs + ".actual"
	err = os.WriteFile(actualPathAbs, actual, 0o600)
	th.AssertNoErr(t, err)

	expected, err := os.ReadFile(fixturePathAbs)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, string(expected), string(actual))
}
