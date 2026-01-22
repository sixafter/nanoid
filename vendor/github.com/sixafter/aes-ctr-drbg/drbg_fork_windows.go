// Copyright (c) 2024-2026 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

//go:build windows

package ctrdrbg

// reseedIfForked is a platform-specific no-op on Windows.
//
// On Windows, the concept of process forking (as in fork(2) on Unix-like systems) does not exist.
// Consequently, the risk of duplicated DRBG state due to process fork is not present, and no
// fork-detection or reseeding logic is required.
//
// This stub function satisfies the interface for fork detection and is used in place of the
// fork-aware implementation found on other platforms. It introduces zero runtime overhead.
//
// Semantics and Behavior:
//   - Does nothing on Windows; safe to call at any point.
//   - Present solely for cross-platform compatibility.
//
// Security Rationale:
//   - There is no process-level forking on Windows.
//   - No action is required to preserve DRBG stream uniqueness in this environment.
//
// Usage Notes:
//   - Automatically included via Go build tags for Windows targets (`//go:build windows`).
//   - All fork-related security logic is compiled out for this platform.
func (d *drbg) reseedIfForked() {
	// No-op: Windows does not implement fork(), so fork detection is unnecessary.
}
