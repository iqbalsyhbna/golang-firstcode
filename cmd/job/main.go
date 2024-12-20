// main.go
package main

import (
	"golang-firstcode/internal/config"
	background "golang-firstcode/internal/job"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: Error loading config: %v", err)
	}

	log.Printf("Loaded config: %+v", config.ApiURL)
	log.Printf("Loaded config: %+v", config.AuthToken)
	log.Printf("Loaded config: %+v", config.DeviceID)
	log.Printf("Loaded config: %+v", config.PlatformID)

	log.Println("Starting background job service...")

	// Start the background job scheduler
	scheduler := background.StartBackgroundJob()

	// Set up a channel to listen for interrupt signals for graceful shutdowns
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for an interrupt signal
	sig := <-sigChan
	log.Printf("Received signal: %s. Shutting down...", sig)

	// Gracefully stop the scheduler
	scheduler.Stop()
	log.Println("Background job service stopped.")
}
