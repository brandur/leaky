package main

import (
)

func main() {
	addLogEntry(LogEntry{term: currentTerm, operation: PUT, data: "foo"})

	go RunState(&server, &clients)
	go RunClient()
	RunServer()
}
