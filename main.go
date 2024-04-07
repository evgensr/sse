package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"
)

var (
	folder = "static"
	//go:embed static/*.html
	content embed.FS
)

func main() {
	r := http.NewServeMux()

	r.HandleFunc("/", pageHandler)
	httpFS, err := fs.Sub(content, ".")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(httpFS))
	r.Handle(fmt.Sprintf("/%s/", folder), fileServer)

	// SSE handler
	r.HandleFunc("/events", eventsHandler)

	log.Println("Listening on port 8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Simulate sending events (you can replace this with real data)
	for i := 0; i < 10; i++ {
		t := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf("Event %s", t))
		time.Sleep(2 * time.Second)
		w.(http.Flusher).Flush()
	}

	// Simulate closing the connection
	closeNotify := r.Context().Done()
	<-closeNotify
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	file, err := content.ReadFile(fmt.Sprintf("%s/page.html", folder))
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(file)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
