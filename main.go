package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	Operation  Operation `json:"operation"`
	Sender     string    `json:"sender"`
	Clock      int8      `json:"clock"`
	DurationMs *int16    `json:"durationMs,omitempty"`
	Peers      *[]string `json:"peers,omitempty"`
	ReturnPeer *string   `json:"returnPeer,omitempty"`
}

type Operation int8

type LocalClock struct {
	mu    sync.Mutex
	value int8
}

const (
	Internal Operation = iota
	ExternalSingle
	ExternalMultiple
	Halt
)

func (c *LocalClock) operationHandler(w http.ResponseWriter, r *http.Request) {
	message := parseRequest(r)
	log.Printf("Got request: %+v", message)

	switch message.Operation {
	case Internal:
		handleInternalOperation(message, c)
	case ExternalSingle:
		handleExternalSingleOperation(message)
	case ExternalMultiple:
		handleExternalMultipleOperation(message)
	case Halt:
		handleHaltOperation(message)
	}

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

func handleInternalOperation(m Message, c *LocalClock) {
	log.Println("Handling internal operation")

	timer := *m.DurationMs
	if timer == 0 {
		timer = 1000
	}

	time.Sleep(time.Duration(timer) * time.Millisecond)

	defer c.mu.Unlock()
	c.mu.Lock()
	c.value++

	log.Println("Internal operation completed, current clock value: ", c.value)
}

func handleExternalSingleOperation(m Message) {
	log.Println("Handling external single operation")
}

func handleExternalMultipleOperation(m Message) {
	log.Println("Handling external multiple operation")
}

func handleHaltOperation(m Message) {
	log.Println("Handling halting operation")
}

func main() {
	mux := http.NewServeMux()

	clock := &LocalClock{value: 0}

	mux.HandleFunc("/", clock.operationHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
