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

func main() {
	log.Println("Starting application...")
    
    // Inicializar aplicación
    application := app.NewConfigApplication()
    log.Println("Application instance created")
    
    // Configurar componentes con manejo de errores
    log.Println("Loading configuration...")
    if err := application.LoadConfig(); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    log.Println("Configuration loaded")
    
    log.Println("Setting up common services...")
    if err := application.SetupCommonServices(); err != nil {
        log.Fatalf("Failed to setup common services: %v", err)
    }
    log.Println("Common services setup complete")
    
    log.Println("Connecting to database...")
    if err := application.SetupDatabase(); err != nil {
        log.Fatalf("Failed to setup database: %v", err)
    }
    log.Println("Database connected")
    defer application.Close()
    
    log.Println("Connecting to Redis...")
    if err := application.SetupRedis(); err != nil {
        log.Fatalf("Failed to setup Redis: %v", err)
    }
    log.Println("Redis connected")

    log.Println("Building router...")
    router := application.BuildRouter()
    log.Println("Router built")
    
    log.Println("Applying security middleware...")
    securedHandler := app.WrapWithSecurityMiddleware(router)
    log.Println("Security middleware applied")

    port := application.GetPort()

    server := &http.Server{
        Addr:           ":" + port,
        Handler:        securedHandler,
        ReadTimeout:    15 * time.Second,
        WriteTimeout:   15 * time.Second,
        IdleTimeout:    60 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    go func() {
        log.Printf("Server started on http://localhost:%s", port)
        log.Println("Press Ctrl+C to stop")
        
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    <-quit
    log.Println("Shutting down server gracefully...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Server forced to shutdown: %v", err)
    }

    log.Println("Server stopped gracefully")
}
