//go:build !windows

package app

import "github.com/shirou/gopsutil/v3/process"

// ShouldTrackApp is a no-op on non-Windows platforms
func ShouldTrackApp(exePath string, proc *process.Process) bool {
	return true
}

// ShouldTrackAppStrict is a no-op on non-Windows platforms
func ShouldTrackAppStrict(exePath string, proc *process.Process) bool {
	return true
}
