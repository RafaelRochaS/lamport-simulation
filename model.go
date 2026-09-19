package main

import "sync"

type Message struct {
	Operation  Operation `json:"operation"`
	Sender     string    `json:"sender"`
	Clock      int       `json:"clock"`
	DurationMs *int16    `json:"durationMs,omitempty"`
	Peers      *[]string `json:"peers,omitempty"`
	ReturnPeer *string   `json:"returnPeer,omitempty"`
}

type LocalClock struct {
	mu    sync.Mutex
	value int
}

type Operation int8

const (
	Internal Operation = iota
	ExternalSingle
	ExternalMultiple
	Halt
)

func (o Operation) String() string {
	switch o {
	case Internal:
		return "Internal"
	case ExternalSingle:
		return "External - Single Call"
	case ExternalMultiple:
		return "External - Multiple Calls"
	case Halt:
		return "Halt"
	default:
		return "Unknown"
	}
}
