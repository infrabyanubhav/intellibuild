package main

import (
	"IntelliBuildCI/v/webhook"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Handle webhook requests
	http.HandleFunc("/webhook", webhook.WebhookHandler)

	// Start the web server
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
