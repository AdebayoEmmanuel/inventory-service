package main

import (
    "log"
    "net/http"
    "os"

    "github.com/AdebayoEmmanuel/inventory-service/internal/handlers"
)

func main() {
    // Get port from environment or use default
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Create a new mux router
    mux := http.NewServeMux()

    // Register handlers
    mux.HandleFunc("/items", handlers.ItemsHandler)
    mux.HandleFunc("/status", handlers.StatusHandler)

    // Start the server
    log.Printf("Server starting on port %s", port)
    if err := http.ListenAndServe(":"+port, mux); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}