package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	Operation  Operation `json:"operation"`
	Sender     string    `json:"sender"`
	Clock      int       `json:"clock"`
	DurationMs *int16    `json:"durationMs,omitempty"`
	Peers      *[]string `json:"peers,omitempty"`
	ReturnPeer *string   `json:"returnPeer,omitempty"`
}

type Operation int8

type LocalClock struct {
	mu    sync.Mutex
	value int
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
		handleInternalOperation(message)
	case ExternalSingle:
		handleExternalSingleOperation(message)
	case ExternalMultiple:
		handleExternalMultipleOperation(message)
	case Halt:
		handleHaltOperation(message)
	}

	increaseClock(c, message.Clock)

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Fatalf("Error writing response: %v", err)
	}
}

func parseRequest(r *http.Request) Message {
	defer r.Body.Close()
	var message Message

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return message
}

func handleInternalOperation(m Message) {
	log.Println("Handling internal operation")

	var timer int16 = 1000
	if m.DurationMs != nil && *m.DurationMs > 0 {
		timer = *m.DurationMs
	}

	log.Println("Timer set to: ", timer)

	time.Sleep(time.Duration(timer) * time.Millisecond)
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

func increaseClock(c *LocalClock, receivedClock int) {
	defer c.mu.Unlock()
	c.mu.Lock()

	maxValue := math.Max(float64(c.value), float64(receivedClock))

	c.value = int(maxValue) + 1

	log.Println("Clock increased, current clock value: ", c.value)
}

func callOperations(c *LocalClock) {
	log.Println("Starting operations caller")

	for {
		sleepTime := rand.Intn(5000)
		log.Println("Operations caller :: sleeping for:", sleepTime)

		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

		log.Println("Operations caller :: woke up, calling random operation")

		message := &Message{
			Clock: c.value,
		}

		handleInternalOperation(*message)
		increaseClock(c, message.Clock)
	}
}

func main() {
	mux := http.NewServeMux()

	clock := &LocalClock{value: 0}
	go callOperations(clock)

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
