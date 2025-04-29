// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package tpm

import (
	"github.com/google/go-tpm/tpm2/transport/windowstpm"
	"tailscale.com/types/opt"
)

func available() opt.Bool {
	t, err := windowstpm.Open()
	if err != nil {
		return opt.NewBool(false)
	}
	defer t.Close()
	return opt.NewBool(true)
}
