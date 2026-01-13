package main

import (
	"godb/web"
	"log"
)

func main() {
	server := web.NewServer(":8080")

	// Initialize database schema
	if err := server.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Start server
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
