package main

import ()

func main() {
	//addLogEntry(LogEntry{term: currentTerm, operation: PUT, data: "foo"})

	peers := newPeers()
	go peers.run()

	server := newServer()

	state := newStateMachine(peers, server)
	go state.run()

	server.run()
}
