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

type Operation int8

type LocalClock struct {
	mu    sync.Mutex
	value int
}
