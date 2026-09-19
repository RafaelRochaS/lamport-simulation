package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func handleInternalOperation(m Message) {
	slog.Info("Handling internal operation")

	var timer int16 = 1000
	if m.DurationMs != nil && *m.DurationMs > 0 {
		timer = *m.DurationMs
	}

	slog.Debug("Timer set to: ", timer)

	time.Sleep(time.Duration(timer) * time.Millisecond)
}

func handleExternalSingleOperation(m Message, clockValue int) {
	slog.Info("Handling external single operation")

	if m.Peers == nil || len(*m.Peers) == 0 {
		slog.Error("No peers provided")
		os.Exit(1)
	}

	peers := *m.Peers

	slog.Debug("External Single Operation :: peers: ", peers)

	peerUrl := fmt.Sprintf("http://%s:8080", peers[0])
	message := Message{
		Operation: Internal,
		Sender:    os.Getenv("PROCESS_ID"),
		Clock:     clockValue,
	}

	slog.Debug("External Single Operation :: calling: ", peerUrl)
	slog.Debug("External Single Operation :: message: ", peerUrl)

	body, err := json.Marshal(message)
	if err != nil {
		slog.Error("Failed to marshal message: %v", err)
		os.Exit(1)
	}

	_, err = http.Post(peerUrl, "application/json", bytes.NewBuffer(body))

	if err != nil {
		slog.Error("Failed to send message: %v", err)
		os.Exit(1)
	}

	slog.Info("External Single Operation :: finished operation call")
}

func handleExternalMultipleOperation(m Message) {
	slog.Info("Handling external multiple operation")
}

func handleHaltOperation(m Message) {
	slog.Info("Handling halting operation")
}
