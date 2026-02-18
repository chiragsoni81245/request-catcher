package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	port := "8080"
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		printRequest(r)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request received\n"))
	})

	http.HandleFunc("/bg-task", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Starting Background Task [waiting 20s...]")
		time.Sleep(20 * time.Second)
		fmt.Printf("Background Task Completed [waiting 20s...]")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("background task completed\n"))
	})

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

