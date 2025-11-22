package app

import "strings"

// DefaultLinux is the default list of user-level process names to ignore on Linux.
// These are typically system or desktop environment processes that are not useful to monitor.
var DefaultLinux = []string{
	// Systemd processes
	"systemd",
	"(sd-pam)",

	// DBus and other daemons
	"dbus-daemon",

	// GNOME services
	".gcr-ssh-agent-",
	".gnome-keyring-",
	".gnome-session-",
	".gnome-shell-wr",
	".at-spi-bus-lau",
	"at-spi2-registryd",
	".gnome-shell-ca",
	"dconf-service",
	".evolution-sour",
	".org.gnome.Shel",
	".evolution-alar",
	".org.gnome.Scre",
	".goa-daemon-wra",
	".goa-identity-s",
	".evolution-cale",
	".evolution-addr",
	"gsd-",   // gnome-settings-daemon prefix
	"gvfsd-", // gnome-virtual-file-system prefix
	"gvfs-",
	"gdm-",

	// XDG and desktop portals
	"xdg-",
	"fusermount3",
	".mutter-x11-fra",
	".localsearch-3-",

	// Pipewire and audio
	"pipewire",
	"pipewire-pulse",
	"wireplumber",
	"speech-dispatcher",

	// Other
	"Xwayland",
	"ssh-agent",
}

// DefaultWindows is the default list of process names to ignore on Windows.
// These are core system processes that are safe to ignore and not relevant for monitoring.
var DefaultWindows = []string{}

// IsIgnored checks if a process name should be ignored based on the ignore list.
// It performs both exact and prefix matching, and handles truncated names that start with a dot.
func IsIgnored(name string, ignoreList []string) bool {
	for _, ignored := range ignoreList {
		if strings.HasSuffix(ignored, "-") {
			// Prefix match (e.g., "gsd-" should match "gsd-color")
			if strings.HasPrefix(name, strings.TrimSuffix(ignored, "-")) {
				return true
			}
			// Also check against the name with a leading dot trimmed, for truncated process names.
			if strings.HasPrefix(strings.TrimPrefix(name, "."), strings.TrimSuffix(ignored, "-")) {
				return true
			}
		} else {
			// Exact match
			if name == ignored {
				return true
			}
		}
	}
	return false
}
