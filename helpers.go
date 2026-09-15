package main

import (
	"encoding/json"
	"log"
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
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return message
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

		operation := rand.Intn(2)

		log.Println("Operations caller :: chosen operation index:", operation)

		if operation == 0 {
			message := &Message{
				Clock: c.value,
			}

			handleInternalOperation(*message)
			increaseClock(c, message.Clock)
		} else if operation == 1 {
			peers := os.Getenv("PEERS")

			if len(peers) < 0 {
				log.Fatalf("No peers found")
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
