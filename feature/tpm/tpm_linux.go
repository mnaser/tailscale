// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package tpm

import (
	"github.com/google/go-tpm/tpm2/transport/linuxtpm"
	"tailscale.com/types/opt"
)

func available() opt.Bool {
	t, err := linuxtpm.Open("/dev/tpm0")
	if err != nil {
		return opt.NewBool(false)
	}
	defer t.Close()
	return opt.NewBool(true)
}
