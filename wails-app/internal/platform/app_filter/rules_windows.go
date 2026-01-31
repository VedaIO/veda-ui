//go:build windows

package app_filter

import (
	"strings"
	"wails-app/internal/platform/executable"
	"wails-app/internal/platform/integrity"

	"github.com/shirou/gopsutil/v3/process"
)

// ShouldExclude returns true if the process is a Windows system component or ProcGuard itself.
func ShouldExclude(exePath string, proc *process.Process) bool {
	exePathLower := strings.ToLower(exePath)

	// Rule 0: Never track ProcGuard itself
	if strings.Contains(exePathLower, "procguard.exe") {
		return true
	}

	// Rule 1: Skip if in System32/SysWOW64 (Windows system processes)
	if strings.Contains(exePathLower, "\\windows\\system32\\") ||
		strings.Contains(exePathLower, "\\windows\\syswow64\\") {
		return true
	}

	// Rule 2: Skip processes with "Microsoft® Windows® Operating System" product name
	productName, err := executable.GetProductName(exePath)
	if err == nil && strings.Contains(productName, "Microsoft® Windows® Operating System") {
		return true
	}

	// Rule 3: Skip system integrity level processes (system services)
	if proc != nil {
		il, err := integrity.GetProcessLevel(uint32(proc.Pid))
		if err == nil && il >= integrity.SystemRID {
			return true
		}
	}

	return false
}
