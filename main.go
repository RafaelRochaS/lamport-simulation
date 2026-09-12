package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Operation  Operation `json:"operation"`
	Sender     string    `json:"sender"`
	Clock      int8      `json:"clock"`
	DurationMs *int8     `json:"durationMs,omitempty"`
	Peers      *[]string `json:"peers,omitempty"`
	ReturnPeer *string   `json:"returnPeer,omitempty"`
}

type Operation int8

const (
	INTERNAL Operation = iota
	EXTERNAL_SINGLE
	EXTERNAL_MULTIPLE
	HALT
)

func operationHandler(w http.ResponseWriter, r *http.Request) {
	message := parseRequest(r)
	log.Printf("Got request: %+v", message)

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Fatalf("Error writing response: %v", err)
	}
}

func parseRequest(r *http.Request) Message {
	var message Message

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return message
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", operationHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
