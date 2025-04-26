// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package ipnlocal

import (
	"tailscale.com/ipn"
	"tailscale.com/tailcfg"
)

// hookSetPeerStatusEnv is the environment passed to the HookSetPeerStatus hook
// to use for setting the status of a peer.
//
// The LocalBackend.mu is locked throughout the call to the hook.
type hookSetPeerStatusEnv struct {
	dead    bool
	b       *LocalBackend // b.mu is locked
	cn      *localNodeContext
	curPeer *tailcfg.NodeView
}

func (e *hookSetPeerStatusEnv) Close() { e.dead = true }

func (e *hookSetPeerStatusEnv) State() ipn.State {
	if e.dead {
		panic("misuse of hook env")
	}
	return e.b.state
}

func (e *hookSetPeerStatusEnv) Peer() tailcfg.NodeView {
	if e.dead {
		panic("misuse of hook env")
	}
	return *e.curPeer
}

func (e *hookSetPeerStatusEnv) PeerHasPeerAPI() bool {
	if e.dead {
		panic("misuse of hook env")
	}
	p := *e.curPeer
	if !p.Valid() {
		return false
	}
	p4, p6 := peerAPIPorts(p)
	return p4 != 0 || p6 != 0
}

func (e *hookSetPeerStatusEnv) PeerHasCap(cap tailcfg.PeerCapability) bool {
	if e.dead {
		panic("misuse of hook env")
	}
	p := *e.curPeer
	if !p.Valid() {
		return false
	}
	if p.Addresses().Len() == 0 {
		return false
	}
	caps := e.cn.PeerCaps(p.Addresses().At(0).Addr())
	// TODO(bradfitz): avoid the PeerCaps map allocs and just ask the boolean
	// question, passing down the cap we want to check. We need to plumb this
	// down into the filter's peercaps method too.
	return caps.HasCapability(cap)
}
