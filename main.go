package main

import (
	"log"
	"net/http"
)

func (c *LocalClock) operationHandler(w http.ResponseWriter, r *http.Request) {
	message := parseRequest(r)
	log.Printf("Got request: %+v", message)

	switch message.Operation {
	case Internal:
		handleInternalOperation(message)
	case ExternalSingle:
		handleExternalSingleOperation(message, c.value)
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
