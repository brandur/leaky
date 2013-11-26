package main

import ()

func main() {
	//addLogEntry(LogEntry{term: currentTerm, operation: PUT, data: "foo"})

	state := newStateMachine(&clients, &server)
	go state.run()

	go RunClient()
	RunServer()
}
