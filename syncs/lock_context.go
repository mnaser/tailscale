// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package syncs

import (
	"context"
	"sync"
)

type mutexKey struct{ *sync.Mutex }

type lockContext struct {
	context.Context
	mu *sync.Mutex
}

// ContextWithLock acquires the given mutex unless it's already held by the parent context.
// It returns a context that tracks the lock state and an unlock function.
// It is a runtime error to pass a nil mutex, to use the returned context after calling unlock,
// or to unlock the parent context before the returned one.
func ContextWithLock(parent context.Context, mu *sync.Mutex) (ctx context.Context, unlock func()) {
	if mu == nil {
		panic("nil mutex")
	}
	c := &lockContext{Context: parent, mu: mu}
	if parent.Value(mutexKey{mu}) != nil {
		return c, c.panicIfParentUnlocked
	}
	mu.Lock()
	return c, c.unlock
}

func (c *lockContext) Value(key any) any {
	if c.mu == nil {
		panic("use of context after unlock")
	}
	if key != (mutexKey{c.mu}) {
		return c.Context.Value(key)
	}
	return c
}

func (c *lockContext) panicIfParentUnlocked() {
	c.Context.Value(mutexKey{c.mu}) // [lockContext.Value] panics if the parent is unlocked
	c.mu = nil
}

func (c *lockContext) unlock() {
	c.mu.Unlock()
	c.mu = nil
}
