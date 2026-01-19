//go:build !darwin
// +build !darwin

// Package app provides stub implementations for platforms where nsync is not available.
// [IMPL:NSYNC_OBSERVER] [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
package app

import (
	"fmt"

	"github.com/fareedst/goful/message"
)

// CopyAll displays a message that nsync is not available on this platform.
// [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func (g *Goful) CopyAll() {
	message.Error(fmt.Errorf("CopyAll (nsync) not available on this platform - use regular copy (c)"))
	g.Copy()
}

// MoveAll displays a message that nsync is not available on this platform.
// [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func (g *Goful) MoveAll() {
	message.Error(fmt.Errorf("MoveAll (nsync) not available on this platform - use regular move (m)"))
	g.Move()
}

// doCopyAll is a stub for platforms without nsync.
// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
func (g *Goful) doCopyAll(sources, destinations []string) {
	message.Error(fmt.Errorf("CopyAll not available on this platform"))
}

// doMoveAll is a stub for platforms without nsync.
// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
func (g *Goful) doMoveAll(sources, destinations []string) {
	message.Error(fmt.Errorf("MoveAll not available on this platform"))
}
