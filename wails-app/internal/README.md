## Internal Package Navigation

The `internal/` package contains all core business logic that is not part of the public API.

### Directory Structure

```
internal/
├── app/                    # Application logic
│   ├── logging.go          # Process monitoring and filtering
│   ├── app_blocklist.go    # App blocking functionality
│   └── screentime.go       # Screen time tracking
├── data/                   # Data layer
│   ├── database.go         # SQLite connection and write queue
│   ├── data_log.go         # Multi-target logging
│   └── config.go           # Configuration management
├── web/                    # Web activity monitoring
│   ├── native-messaging.go # Native messaging host
│   ├── extension.go        # Extension detection
│   └── web_blocklist.go    # Website blocking
├── daemon/                 # Background services
│   └── daemon.go           # Service orchestration
└── platform/               # OS-specific implementations
    ├── nativehost/         # Browser integration
    ├── autostart/          # Windows autostart
    ├── icon/               # App icon extraction
    └── [other platforms]/
```

### Module Roles

#### 1. **app/** - Core Application Logic
- **logging.go**: Process event detection, filtering rules, database writes
- **app_blocklist.go**: Manages blocked applications list
- **screentime.go**: Tracks foreground window time with buffering

#### 2. **data/** - Persistence Layer
- **database.go**: SQLite initialization, async write queue pattern
- **data_log.go**: Dual logging (file + database) with singleton pattern
- **config.go**: JSON-based configuration storage

#### 3. **web/** - Browser Integration
- **native-messaging.go**: Chrome Native Messaging protocol implementation
- **extension.go**: Extension installation detection
- **web_blocklist.go**: Website blocking and synchronization

#### 4. **daemon/** - Background Services
- **daemon.go**: Launches three monitoring goroutines:
  - Process event logger (2s interval)
  - Screen time tracker (1s interval)
  - App blocklist enforcer (100ms loop)

#### 5. **platform/** - OS Abstractions
Platform-specific implementations using Go build tags:
- `*_windows.go`: Windows-specific code
- `*_darwin.go`: macOS stubs (no-op)
- `*_linux.go`: Linux stubs (no-op)

### Key Patterns

1. **Async Write Queue**: Single database writer goroutine prevents lock contention
2. **Multi-Target Logging**: All logs go to both file and database simultaneously
3. **Platform Abstraction**: Clean separation of OS-specific code
4. **Process State Management**: Daemon maintains state for process detection and screen time

### Data Flow

```
Frontend → API Layer → Internal Logic → Async Queue → SQLite Database
    ↓           ↓              ↓              ↓              ↓
Extension ← Native Host ← Web Package ← Async Queue ← SQLite Database
```

---

## Module Dependencies

### Core Dependencies
- All modules depend on `data/` for database access
- `daemon/` orchestrates services from `app/` and `web/`
- `platform/` provides OS-specific implementations to other modules

### External Dependencies
- **golang.org/x/sys/windows**: Windows API access
- **github.com/shirou/gopsutil**: Process enumeration
- **github.com/wailsapp/wails/v2**: Application framework
