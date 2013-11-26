package main

import ()

func main() {
	//addLogEntry(LogEntry{term: currentTerm, operation: PUT, data: "foo"})

	server := newServer()

	state := newStateMachine(&clients, server)
	go state.run()

	go RunClient()
	server.run()
}
