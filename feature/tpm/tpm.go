// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

// Package tpm implements support for TPM 2.0 devices.
package tpm

import (
	"sync"

	"tailscale.com/feature"
	"tailscale.com/hostinfo"
	"tailscale.com/tailcfg"
)

var availableOnce = sync.OnceValue(available)

func init() {
	feature.Register("tpm")
	hostinfo.RegisterHostinfoNewHook(func(hi *tailcfg.Hostinfo) {
		hi.TPMAvailable = availableOnce()
	})
}
