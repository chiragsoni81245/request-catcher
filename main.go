package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func withMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Server-Name", os.Getenv("HOSTNAME"))
		handler(w, r)
	}
}

func main() {
	port := "8080"

	http.HandleFunc("/", withMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	}))

	http.HandleFunc("/check", withMiddleware(func(w http.ResponseWriter, r *http.Request) {
		printRequest(r)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request received\n"))
	}))

	http.HandleFunc("/bg-task", withMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Starting Background Task [waiting 20s...] \n")

		select {
		case <-time.After(20 * time.Second):
			fmt.Printf("Background Task Completed [waiting 20s...] \n")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("background task completed\n"))
		case <-r.Context().Done():
			fmt.Printf("Client aborted request\n")
			return
		}
	}))

	http.HandleFunc("/session-check/", withMiddleware(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/session-check/")
		if name == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "name is required"})
			return
		}

		filePath := filepath.Join("sessions", name)
		_, err := os.Stat(filePath)
		exists := !os.IsNotExist(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":   name,
			"exists": exists,
		})
	}))

	http.HandleFunc("/session-create/", withMiddleware(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/session-create/")
		if name == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "name is required"})
			return
		}

		// Ensure sessions directory exists
		if err := os.MkdirAll("sessions", 0755); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
			return
		}

		filePath := filepath.Join("sessions", name)
		file, err := os.Create(filePath)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
			return
		}
		file.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":    name,
			"created": true,
		})
	}))

	log.Printf("Starting server on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func printRequest(r *http.Request) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Time: %s\n", time.Now().Format(time.RFC3339))
	fmt.Println(strings.Repeat("-", 80))

	// Request line
	fmt.Printf("Method:  %s\n", r.Method)
	fmt.Printf("URL:     %s\n", r.URL.String())
	fmt.Printf("Proto:   %s\n", r.Proto)
	fmt.Printf("Host:    %s\n", r.Host)
	fmt.Printf("Remote:  %s\n", r.RemoteAddr)
	fmt.Printf("Length:  %d\n", r.ContentLength)

	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("Headers:")

	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", name, value)
		}
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("Body:")

	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("  Error reading body: %v\n", err)
		} else if len(bodyBytes) > 0 {
			fmt.Println(indent(string(bodyBytes), "  "))
		} else {
			fmt.Println("  <empty>")
		}
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
}

func indent(text, prefix string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}
