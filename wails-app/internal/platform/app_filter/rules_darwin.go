//go:build darwin

package app_filter

import (
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// ShouldExclude returns true if the process should be ignored (Stub for non-Windows).
func ShouldExclude(exePath string, proc *process.Process) bool {
	exePathLower := strings.ToLower(exePath)

	// Rule 0: Never track ProcGuard itself
	if strings.Contains(exePathLower, "procguard") {
		return true
	}

	// Basic filtering for Linux/Darwin could be added here (e.g., /usr/bin vs /bin)
	return false
}
