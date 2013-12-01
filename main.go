package main

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
)

func generateName() Name {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	name := Name("leaky-" + id.String())
	fmt.Printf("name=%v\n", name)
	return name
}

func main() {
	//addLogEntry(LogEntry{term: currentTerm, operation: PUT, data: "foo"})
	name := generateName()

	peerSet := newPeerSet()
	go peerSet.run()

	// point to an existing cluster if one was given
	if conf.peerUrl != "" {
		peerSet.addPeer(conf.peerUrl)
		peerSet.JoinRequestChan <- JoinRequest{Name: name, Url: conf.http}
	}
	// @todo: else self-join immediately

	server := newServer()
	go server.run()

	state := newStateMachine(name, peerSet, server)
	state.run()
}
