package main

import ()

func main() {
	//addLogEntry(LogEntry{term: currentTerm, operation: PUT, data: "foo"})

	peerSet := newPeerSet()
	go peerSet.run()

	server := newServer()

	state := newStateMachine(peerSet, server)
	go state.run()

	server.run()
}
