package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	renderURL := os.Getenv("RENDER_URL")
	if renderURL == "" {
		log.Fatal("RENDER_URL environment variable is not set")
	}

	go keepAlive(renderURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Keep-alive Server is Running")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func keepAlive(url string) {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error making request to %s: %v", url, err)
			continue
		}
		resp.Body.Close()
		log.Printf("Ping sent to %s, status: %s", url, resp.Status)
	}
}
