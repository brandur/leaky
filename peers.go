package main

import (
	"encoding/json"
)

type Peers struct {
	set []Client
	AppendEntriesRequestChan  chan AppendEntriesRequest
	AppendEntriesResponseChan chan AppendEntriesResponse
	RequestVoteRequestChan    chan RequestVoteRequest
	RequestVoteResponseChan   chan RequestVoteResponse
}


func newPeers() *Peers {
	return &Peers{
		AppendEntriesRequestChan:  make(chan AppendEntriesRequest, 50),
		AppendEntriesResponseChan: make(chan AppendEntriesResponse, 50),
		RequestVoteRequestChan:    make(chan RequestVoteRequest, 50),
		RequestVoteResponseChan:   make(chan RequestVoteResponse, 50),
	}
}

func (p *Peers) broadcast(path string, data []byte, handler func([]byte)) {
	for i := range p.set {
		go p.set[i].connectAndSend(path, data, handler)
	}
}

func (p *Peers) run() {
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

func (p *Peers) buildHandler(handler func(*Peers, []byte)) func([]byte) {
	return func(responseData []byte) {
		handler(p, responseData)
	}
}

func handleAppendEntriesResponse(p *Peers, responseData []byte) {
	response := AppendEntriesResponse{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		panic(err)
	}
	p.AppendEntriesResponseChan <- response
}

func handleRequestVoteResponse(p *Peers, responseData []byte) {
	response := RequestVoteResponse{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		panic(err)
	}
	p.RequestVoteResponseChan <- response
}
