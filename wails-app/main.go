package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"log"
	"os"
	"strings"
	"wails-app/api"
	"wails-app/internal/daemon"
	"wails-app/internal/data"
	"wails-app/internal/web"
)

// Embed the entire frontend/dist directory into the Go binary
// This allows the app to be distributed as a single executable
//
//go:embed all:frontend/dist
var assets embed.FS

// startup is called when the Wails app starts
// The context is saved so we can call runtime methods (WindowShow, etc.) later
//
// Responsibilities:
//  1. Save the Wails runtime context for later use
//  2. Initialize the database
//  3. Initialize the logger
//  4. Create the API server
//  5. Start the background daemon for process/web monitoring
//  6. Start the native messaging host for Chrome extension communication
func (a *App) startup(ctx context.Context) {
	// Save context - CRITICAL for calling ShowWindow() and other runtime methods
	a.ctx = ctx

	// Initialize database connection
	db, err := data.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize logger with database
	data.NewLogger(db)

	// Create API server with database connection
	a.Server = api.NewServer(db)

	// Start the background daemon that monitors processes and web activity
	// This runs independently of the GUI - continues even when window is hidden
	daemon.StartDaemon(a.Server.Logger, db)

	// Ensure Native Messaging Host is registered
	// This creates the registry key and manifest file so Chrome can find us
	// We do this on every startup to ensure the config is correct
	if err := web.RegisterExtension("hkanepohpflociaodcicmmfbdaohpceo"); err != nil {
		log.Printf("Failed to register Store extension: %v", err)
	}
	if err := web.RegisterExtension("gpaafgcbiejjpfdgmjglehboafdicdjb"); err != nil {
		log.Printf("Failed to register Dev extension: %v", err)
	}
}

func main() {
	// CRITICAL: Detect if we're launched as a native messaging host
	// Native messaging hosts have Stdin connected to a pipe (Chrome's stdout)
	// GUI launches have Stdin as a terminal/invalid
	stat, _ := os.Stdin.Stat()
	isNativeMessagingMode := (stat.Mode() & os.ModeCharDevice) == 0
	
	if isNativeMessagingMode {
		// We're a native messaging host - run messaging loop only, no GUI
		log.Println("Starting in native messaging mode (no GUI)")
		web.Run()
		// When web.Run() exits (Chrome disconnects), this process terminates
		return
	}
	
	// Normal GUI mode - create app instance
	log.Println("Starting in GUI mode")
	app := NewApp()

	// Check if this is a native messaging launch (first instance)
	// Chrome passes the extension origin as an argument: chrome-extension://<id>/
	isNativeMessaging := false
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "chrome-extension://") {
			isNativeMessaging = true
			break
		}
	}
	
	// Fallback to WD check
	if !isNativeMessaging {
		wd, _ := os.Getwd()
		if wd != "" {
			wdLower := strings.ToLower(wd)
			if strings.Contains(wdLower, "chrome") || strings.Contains(wdLower, "google") {
				isNativeMessaging = true
			}
		}
	}

	if isNativeMessaging {
		log.Println("First instance launched by Chrome (Native Messaging)")
		app.IsNativeMessagingActive = true
	}

	// Create and run the Wails application with configuration
	err := wails.Run(&options.App{
		Title:  "ProcGuard",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,

		// HideWindowOnClose: Hide window instead of closing application when X is clicked
		// Why: Allows daemon to keep running in background while window is hidden
		// User can reopen by double-clicking executable (SingleInstanceLock handles this)
		HideWindowOnClose: true,

		// SingleInstanceLock: Prevent multiple instances of the application
		//
		// Without this, running the executable multiple times creates multiple processes.
		// Combined with HideWindowOnClose=true, this would cause hidden processes to accumulate.
		//
		// With SingleInstanceLock, only one GUI process can run at a time.
		// Native messaging instances bypass this entirely (they exit main() early).
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "com.procguard.wails-app",
			OnSecondInstanceLaunch: func(data options.SecondInstanceData) {
				// This only fires for GUI launches (native messaging instances already exited)
				log.Println("Second GUI instance detected - showing existing window")
				app.ShowWindow()
			},
		},
		
		// Bind the app struct to make its methods available to frontend JS
		// Frontend can call these via window.go.main.App.MethodName()
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
