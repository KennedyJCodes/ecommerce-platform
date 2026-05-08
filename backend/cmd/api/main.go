// The main package is the entry point for the sales application.
// Configures and initializes all necessary services, including the database,
// Redis, HTTP router, and signal handling for elegant shutdown.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/app"
	_ "github.com/go-sql-driver/mysql"
)

// main initializes and runs the web application.

// Performs the following operations in order:
//  1. Creates the application instance
//  2. Loads the configuration from environment variables or files
//  3. Configures common services (logging, validation, etc.)
//  4. Establishes a connection to the MySQL database
//  5. Establishes a connection to Redis for caching
//  6. Builds the HTTP router with all routes
//  7. Applies security middleware
//  8. Starts the HTTP server
//  9. Waits for an interrupt signal for a graceful shutdown

// The function terminates program execution if any initialization step fails.
// The server gracefully stops upon receiving a SIGINT or SIGTERM,
// allowing up to 30 seconds to complete ongoing requests.

func main() {
	log.Println("Starting application...")
    
    // Create application instance.
    application := app.NewConfigApplication()
    log.Println("Application instance created")
    
    // Load configuration from environment variables or config files
    log.Println("Loading configuration...")
    if err := application.LoadConfig(); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    log.Println("Configuration loaded")
    
    // Configure common services such as validators, loggers, etc.
    log.Println("Setting up common services...")
    if err := application.SetupCommonServices(); err != nil {
        log.Fatalf("Failed to setup common services: %v", err)
    }
    log.Println("Common services setup complete")
    
    // Establish a connection to a MySQL database
    log.Println("Connecting to database...")
    if err := application.SetupDatabase(); err != nil {
        log.Fatalf("Failed to setup database: %v", err)
    }
    log.Println("Database connected")
    defer application.Close()
    
    // Establish a connection with Redis for caching and sessions
    log.Println("Connecting to Redis...")
    if err := application.SetupRedis(); err != nil {
        log.Fatalf("Failed to setup Redis: %v", err)
    }
    log.Println("Redis connected")

    // Build router with all routes and handlers
    log.Println("Building router...")
    router := application.BuildRouter()
    log.Println("Router built")
    
    // Apply security middleware (CORS, rate limiting, etc.)v
    log.Println("Applying security middleware...")
    securedHandler := app.WrapWithSecurityMiddleware(router)
    log.Println("Security middleware applied")

    port := application.GetPort()

    // Configure HTTP server with appropriate timeouts
    server := &http.Server{
        Addr:           ":" + port,
        Handler:        securedHandler,
        ReadTimeout:    15 * time.Second,
        WriteTimeout:   15 * time.Second,
        IdleTimeout:    60 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    // Channel for receiving interruption signals
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    // Start server in separate goroutine
    go func() {
        log.Printf("Server started on http://localhost:%s", port)
        log.Println("Press Ctrl+C to stop")
        
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Wait for interrupt signal
    <-quit
    log.Println("Shutting down server gracefully...")

    // Create context with 30-second timeout for graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Attempt to shut down the server gracefully
    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Server forced to shutdown: %v", err)
    }

    log.Println("Server stopped gracefully")
}
