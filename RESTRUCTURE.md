# ProcGuard Codebase Structure

## Design Principle

**Feature-based platform packages** - Each platform feature is in its own package with OS-specific implementations side by side. No interfaces, no abstraction layers, just direct imports.

## Current Structure

```
wails-app/internal/
├── app/                           # Application monitoring domain
│   ├── screentime.go              # Uses platform/screentime directly
│   ├── filter.go                  # App filtering (Windows-only)
│   ├── process.go                 # Process monitoring
│   ├── icon.go                    # Icon extraction (to migrate)
│   └── integrity.go               # Integrity levels (to migrate)
│
├── platform/                      # OS-specific implementations
│   └── screentime/                # Foreground window tracking
│       ├── screentime_windows.go  # CGO implementation
│       └── screentime_darwin.go   # macOS stub
│
├── data/                          # Database layer
├── daemon/                        # Daemon orchestration
└── web/                           # Chrome extension
```

## How It Works

Callers import directly:
```go
import "wails-app/internal/platform/screentime"

info := screentime.GetActiveWindowInfo()  // Direct call!
```

Build tags (`//go:build windows` / `//go:build darwin`) handle OS selection automatically.

## Future Migrations

When ready, move these to `platform/`:

| Current | Target |
|---------|--------|
| `app/icon.go` | `platform/icon/icon_windows.go` |
| `app/integrity.go` | `platform/integrity/integrity_windows.go` |
| Window enumeration (in process.go) | `platform/window/window_windows.go` |
| Publisher extraction (in process.go) | `platform/executable/executable_windows.go` |

Each will have a corresponding `_darwin.go` stub.
