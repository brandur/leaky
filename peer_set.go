package main

import (
	"encoding/json"
)

type PeerSet struct {
	set []Client
	AppendEntriesRequestChan  chan AppendEntriesRequest
	AppendEntriesResponseChan chan AppendEntriesResponse
	RequestVoteRequestChan    chan RequestVoteRequest
	RequestVoteResponseChan   chan RequestVoteResponse
}


func newPeerSet() *PeerSet {
	return &PeerSet{
		AppendEntriesRequestChan:  make(chan AppendEntriesRequest, 50),
		AppendEntriesResponseChan: make(chan AppendEntriesResponse, 50),
		RequestVoteRequestChan:    make(chan RequestVoteRequest, 50),
		RequestVoteResponseChan:   make(chan RequestVoteResponse, 50),
	}
}

func (p *PeerSet) broadcast(path string, data []byte, handler func([]byte)) {
	for i := range p.set {
		go p.set[i].connectAndSend(path, data, handler)
	}
}

func (p *PeerSet) buildHandler(handler func(*PeerSet, []byte)) func([]byte) {
	return func(responseData []byte) {
		handler(p, responseData)
	}
}

func handleAppendEntriesResponse(p *PeerSet, responseData []byte) {
	response := AppendEntriesResponse{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		panic(err)
	}
	p.AppendEntriesResponseChan <- response
}

func handleRequestVoteResponse(p *PeerSet, responseData []byte) {
	response := RequestVoteResponse{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		panic(err)
	}
	p.RequestVoteResponseChan <- response
}

func (p *PeerSet) len() int {
	return len(p.set)
}

func (p *PeerSet) run() {
	for {
		select {
		case request := <-p.AppendEntriesRequestChan:
			requestData, err := json.Marshal(request)
			if err != nil {
				panic(err)
			}
			p.broadcast("/append-entries", requestData, p.buildHandler(handleAppendEntriesResponse))

		case request := <-p.RequestVoteRequestChan:
			requestData, err := json.Marshal(request)
			if err != nil {
				panic(err)
			}
			p.broadcast("/request-vote", requestData, p.buildHandler(handleRequestVoteResponse))
		}
	}
}
