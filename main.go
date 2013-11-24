package main

import (
)

func main() {
	addLogEntry(LogEntry{term: 1, operation: PUT, data: "foo"})

	go RunState(&server)
	RunServer()
}
