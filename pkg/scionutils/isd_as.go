// Copyright 2026 OVGU Magdeburg
// SPDX-License-Identifier: Apache-2.0

package scionutils

import (
	"strings"

	"github.com/scionproto/scion/pkg/addr"
)

func GetISDFromISDAS(isdAS string) string {
	parts := strings.Split(isdAS, "-")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

func IsValidISDAS(isdAs string) bool {
	_, err := addr.ParseIA(isdAs)
	return err == nil
}
