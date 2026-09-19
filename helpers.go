package main

import (
	"encoding/json"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func parseRequest(r *http.Request) Message {
	defer r.Body.Close()
	var message Message

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	if err != nil {
		slog.Error("Error parsing JSON: %v", err)
		os.Exit(1)
	}

	return message
}

func increaseClock(c *LocalClock, receivedClock int) {
	defer c.mu.Unlock()
	c.mu.Lock()

	maxValue := math.Max(float64(c.value), float64(receivedClock))

	c.value = int(maxValue) + 1

	slog.Info("Clock increased, current clock value: ", c.value)
}

func callOperations(c *LocalClock) {
	slog.Info("Starting operations caller")

	for {
		sleepTime := rand.Intn(5000)
		slog.Debug("Operations caller :: sleeping for:", sleepTime)

		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

		slog.Debug("Operations caller :: woke up, calling random operation")

		operation := rand.Intn(2)

		slog.Debug("Operations caller :: chosen operation index: %s", operation)

		if operation == 0 {
			message := &Message{
				Clock: c.value,
			}

			handleInternalOperation(*message)
			increaseClock(c, message.Clock)
		} else if operation == 1 {
			peers := os.Getenv("PEERS")

			if len(peers) < 0 {
				slog.Error("No peers found")
				os.Exit(1)
			}

			peersList := strings.Split(peers, ",")
			peerValue := rand.Intn(4)
			var peerChosen string

			if len(peersList) < peerValue {
				peerChosen = peersList[0]
			} else {
				peerChosen = peersList[peerValue]
			}

			peersFinal := []string{peerChosen}

			message := &Message{
				Operation: ExternalSingle,
				Sender:    os.Getenv("PROCESS_ID"),
				Clock:     c.value,
				Peers:     &peersFinal,
			}

			handleExternalSingleOperation(*message, c.value)
			increaseClock(c, message.Clock)
		}
	}
}
