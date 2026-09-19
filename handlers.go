package main

import (
	"log/slog"
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
	message := Message{
		Operation: Internal,
		Sender:    os.Getenv("PROCESS_ID"),
		Clock:     clockValue,
	}
	callPeer(peers[0], &message)

	slog.Debug("External Single Operation :: peers: ", peers)

	slog.Info("External Single Operation :: finished operation call")
}

func handleExternalMultipleOperation(m Message, clockValue int) {
	slog.Info("Handling external multiple operation")

	if m.Peers == nil || len(*m.Peers) == 0 {
		slog.Error("No peers provided")
		os.Exit(1)
	}

	peers := *m.Peers
	shiftedArray := peers[1:]

	var operation Operation
	if len(shiftedArray) > 1 {
		operation = ExternalMultiple
	} else {
		operation = ExternalSingle
	}

	message := Message{
		Operation:  operation,
		Sender:     os.Getenv("PROCESS_ID"),
		Clock:      clockValue,
		Peers:      &shiftedArray,
		ReturnPeer: m.ReturnPeer,
	}

	callPeer(shiftedArray[0], &message)
}

func handleHaltOperation(m Message) {
	slog.Info("Handling halting operation")
}
