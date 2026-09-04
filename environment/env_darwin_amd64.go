// Copyright 2026 OVGU Magdeburg
// SPDX-License-Identifier: Apache-2.0

package environment

import "embed"

func init() {
	EndhostEnv = &EndhostEnvironment{
		BasePath:          "/Applications/scion/",
		ConfigPath:        "/Applications/scion/",
		DaemonConfigPath:  "/Applications/scion/",
		ControlConfigPath: "/Applications/scion/",
		RouterConfigPath:  "/Applications/scion/",
		DatabasePath:      "/Applications/scion/database/",
		LogPath:           "/Applications/scion/logs/",
	}
}
