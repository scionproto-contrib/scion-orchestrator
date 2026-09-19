// Copyright 2026 OVGU Magdeburg
// SPDX-License-Identifier: Apache-2.0

package environment

func init() {
	HostEnv = &HostEnvironment{
		BasePath:          "/Applications/scion/",
		ConfigPath:        "/Applications/scion/",
		DaemonConfigPath:  "/Applications/scion/",
		ControlConfigPath: "/Applications/scion/",
		RouterConfigPath:  "/Applications/scion/",
		DatabasePath:      "/Applications/scion/database/",
		LogPath:           "/Applications/scion/logs/",
	}
}
