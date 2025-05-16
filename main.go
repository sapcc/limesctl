// SPDX-FileCopyrightText: 2018 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"time"

	"github.com/sapcc/go-api-declarations/bininfo"
	"github.com/sapcc/go-bits/httpext"

	"github.com/sapcc/limesctl/v3/internal/cmd"
)

func main() {
	// using a short timeout is acceptable here since this process is not a server
	ctx := httpext.ContextWithSIGINT(context.Background(), 100*time.Millisecond)

	v := &cmd.VersionInfo{
		Version:       bininfo.VersionOr("dev"),
		GitCommitHash: bininfo.CommitOr("unknown"),
		BuildDate:     bininfo.BuildDateOr("now"),
	}
	cmd.Execute(ctx, v)
}
