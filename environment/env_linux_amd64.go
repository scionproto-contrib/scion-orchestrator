// Copyright 2026 OVGU Magdeburg
// SPDX-License-Identifier: Apache-2.0

package environment

func init() {
	HostEnv = &HostEnvironment{
		BasePath:          "/etc/scion/",
		ConfigPath:        "/etc/scion/",
		DaemonConfigPath:  "/etc/scion/",
		ControlConfigPath: "/etc/scion/",
		RouterConfigPath:  "/etc/scion/",
		DatabasePath:      "/var/lib/scion/",
		LogPath:           "/var/log/scion/",
	}
}
