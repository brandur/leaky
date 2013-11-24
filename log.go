package main

type LogOperation int

const (
	PUT    LogOperation = 0
	DELETE LogOperation = 1
)

type LogEntry struct {
	term      int
	operation LogOperation
	data      string
}

var (
	log   []LogEntry
	index int = 0
)

func init() {
	log = make([]LogEntry, 1000)
}

func addLogEntry(logEntry LogEntry) {
	log[index] = logEntry
}
