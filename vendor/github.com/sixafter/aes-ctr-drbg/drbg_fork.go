// Copyright (c) 2024-2026 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

//go:build !windows

package ctrdrbg

import (
	"os"
	"sync/atomic"
)

// reseedIfForked performs fork detection and conditional reseeding according to the configured interval.
//
// This method ensures the DRBG instance reseeds its state in the event of a process fork (e.g., fork(2)).
// The fork detection check is performed on every output request if ForkDetectionInterval is zero (the default, safest setting).
// If ForkDetectionInterval is set to a nonzero value N, the fork detection check is only performed every N output requests.
//
// Semantics and Behavior:
//   - Increments the output request counter atomically.
//   - If ForkDetectionInterval is zero, always checks for fork (fully compliant, safest).
//   - If ForkDetectionInterval is N>0, checks only every Nth output request (performance-tuned, non-compliant).
//   - If a fork is detected (current PID != cached PID), reseeds the DRBG instance and updates the cached PID.
//
// Security Rationale:
//   - Forking a process duplicates DRBG state. If not reseeded, parent and child may produce identical output, violating forward secrecy.
//   - Interval-based fork detection reduces syscall overhead but introduces a risk window; do not use in regulated/compliance-critical environments.
//
// Usage:
//   - Call at the beginning of every output-generating operation (e.g., Read, ReadWithAdditionalInput).
func (d *drbg) reseedIfForked() {
	interval := d.config.ForkDetectionInterval
	if interval == 0 {
		// Default: always check for fork
		current := os.Getpid()
		if current != d.pid {
			_ = d.Reseed(nil) // Best-effort reseed
			d.pid = current
		}
		return
	}

	// Only check every Nth request
	n := atomic.AddUint64(&d.requests, 1)
	if n%interval != 0 {
		return // Not time to check yet
	}
	current := os.Getpid()
	if current != d.pid {
		_ = d.Reseed(nil) // Best-effort reseed
		d.pid = current
	}
}
