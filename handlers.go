package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func handleInternalOperation(m Message) {
	log.Println("Handling internal operation")

	var timer int16 = 1000
	if m.DurationMs != nil && *m.DurationMs > 0 {
		timer = *m.DurationMs
	}

	log.Println("Timer set to: ", timer)

	time.Sleep(time.Duration(timer) * time.Millisecond)
}

func handleExternalSingleOperation(m Message, clockValue int) {
	log.Println("Handling external single operation")

	if m.Peers == nil || len(*m.Peers) == 0 {
		log.Fatalf("No peers provided")
	}

	peers := *m.Peers
	peerUrl := fmt.Sprintf("http://%s:8080", peers[0])
	message := Message{
		Operation: Internal,
		Sender:    os.Getenv("PROCESS_ID"),
		Clock:     clockValue,
	}

	body, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	_, err = http.Post(peerUrl, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}

func handleExternalMultipleOperation(m Message) {
	log.Println("Handling external multiple operation")
}

func handleHaltOperation(m Message) {
	log.Println("Handling halting operation")
}
