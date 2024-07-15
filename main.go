package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Render Keep-Alive Server</title>
</head>
<body>
    <h1>Keep-Alive Server is Running</h1>
    <p>Last ping sent: %s</p>
</body>
</html>
`

var lastPingTime string

func main() {
	renderURL := os.Getenv("RENDER_URL")
	if renderURL == "" {
		log.Fatal("RENDER_URL environment variable is not set")
	}

	go keepAlive(renderURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, htmlTemplate, lastPingTime)
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
	for ; true; <-ticker.C {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error making request to %s: %v", url, err)
			continue
		}
		resp.Body.Close()
		lastPingTime = time.Now().Format(time.RFC1123)
		log.Printf("Ping sent to %s, status: %s", url, resp.Status)
	}
}
