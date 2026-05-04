// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/sapcc/go-api-declarations/limes"
)

const (
	csvHeaderClusterID   = "cluster id"
	csvHeaderDomainID    = "domain id"
	csvHeaderDomainName  = "domain name"
	csvHeaderProjectID   = "project id"
	csvHeaderProjectName = "project name"

	csvHeaderArea     = "area"
	csvHeaderService  = "service"
	csvHeaderCategory = "category"
	csvHeaderResource = "resource"
	csvHeaderRate     = "rate"

	csvHeaderCapacity      = "capacity"
	csvHeaderQuota         = "quota"
	csvHeaderProjectsQuota = "projects quota"
	csvHeaderDomainsQuota  = "domains quota"
	csvHeaderUsage         = "usage"
	csvHeaderPhysicalUsage = "physical usage"
	csvHeaderLimit         = "limit"
	csvHeaderDefaultLimit  = "default limit"
	csvHeaderWindow        = "window"
	csvHeaderDefaultWindow = "default window"
	csvHeaderUnit          = "unit"
	csvHeaderScrapedAt     = "scraped at (UTC)"
)

func timestampToString(timestamp *limes.UnixEncodedTime) string {
	if timestamp == nil {
		return ""
	}
	return timestamp.Format(time.RFC3339)
}

func zeroIfNil(ptr *uint64) uint64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func emptyStrIfNil(ptr *uint64, formatter ValueFormatter) string {
	if ptr == nil {
		return ""
	}
	return formatter(*ptr)
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
