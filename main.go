package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func (c *LocalClock) operationHandler(w http.ResponseWriter, r *http.Request) {
	message := parseRequest(r)
	slog.Debug("Got request: %+v", message)

	switch message.Operation {
	case Internal:
		increaseClock(c, message.Clock)
		handleInternalOperation(message)
	case ExternalSingle:
		increaseClock(c, message.Clock)
		handleExternalSingleOperation(message, c.value)
	case ExternalMultiple:
		increaseClock(c, message.Clock)
		handleExternalMultipleOperation(message)
	case Halt:
		handleHaltOperation(message)
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		slog.Error("Error writing response: %v", err)
		os.Exit(1)
	}
}

func setUpLogging() {
	logFile, err := os.OpenFile("/app/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
		os.Exit(1)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	h := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: getLogLevel()})
	slog.SetDefault(slog.New(h))
}

func getLogLevel() slog.Level {
	envVar := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(envVar) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func main() {
	setUpLogging()

	mux := http.NewServeMux()

	clock := &LocalClock{value: 0}
	go callOperations(clock)

	mux.HandleFunc("/", clock.operationHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	slog.Info("Starting server on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Error starting server: ", err)
		os.Exit(1)
	}
}
