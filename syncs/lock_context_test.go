// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package syncs

import (
	"context"
	"sync"
	"testing"
)

func TestLockContext_HappyPath(t *testing.T) {
	var mu sync.Mutex
	parent, unlockParent := ContextWithLock(context.Background(), &mu)
	wantLocked(t, &mu) // mu is locked by parent

	child, unlockChild := ContextWithLock(parent, &mu)
	wantLocked(t, &mu) // mu is still locked by parent

	var mu2 sync.Mutex
	context2, unlock2 := ContextWithLock(child, &mu2)
	_ = context2
	wantLocked(t, &mu2)   // mu2 is locked by context2
	unlock2()             // unlocks mu2
	wantUnlocked(t, &mu2) // mu2 is now unlocked

	unlockChild()      // noop
	wantLocked(t, &mu) // mu is still locked by parent

	unlockParent()       // unlocks mu
	wantUnlocked(t, &mu) // mu is now unlocked
}

func TestLockContext_UnlockParentFirst(t *testing.T) {
	var mu sync.Mutex
	parent, unlockParent := ContextWithLock(context.Background(), &mu)
	_, unlockChild := ContextWithLock(parent, &mu)

	unlockParent()            // unlocks mu
	wantUnlocked(t, &mu)      // mu is now unlocked
	wantPanic(t, unlockChild) // panics since mu is already unlocked by parent
}

func wantPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		recover()
	}()
	fn()
	t.Fatal("failed to panic")
}

func wantLocked(t *testing.T, m *sync.Mutex) {
	if m.TryLock() {
		m.Unlock()
		t.Fatal("mutex is not locked")
	}
}

func wantUnlocked(t *testing.T, m *sync.Mutex) {
	t.Helper()
	if !m.TryLock() {
		t.Fatal("mutex is locked")
	}
	m.Unlock()
}
