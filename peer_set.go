package main

type PeerSet struct {
	AppendEntriesRequestChan  chan AppendEntriesRequest
	AppendEntriesResponseChan chan AppendEntriesResponse
	RequestVoteRequestChan    chan RequestVoteRequest
	RequestVoteResponseChan   chan RequestVoteResponse
	set                       []*Client
}

func newPeerSet() *PeerSet {
	return &PeerSet{
		AppendEntriesRequestChan:  make(chan AppendEntriesRequest, 50),
		AppendEntriesResponseChan: make(chan AppendEntriesResponse, 50),
		RequestVoteRequestChan:    make(chan RequestVoteRequest, 50),
		RequestVoteResponseChan:   make(chan RequestVoteResponse, 50),
		set:                       make([]*Client, 0, 100),
	}
}

func (p *PeerSet) addClient(url string) {
	// grow into cap by one element and insert the new client there
	p.set = p.set[:len(p.set)+1]
	p.set[len(p.set)] = &Client{
		PeerSet: p,
		Url: url,
	}
}

func (p *PeerSet) len() int {
	return len(p.set)
}

func (p *PeerSet) run() {
	for {
		select {
		case request := <-p.AppendEntriesRequestChan:
			for i := range p.set {
				go p.set[i].sendAppendEntriesRequest(&request)
			}

		case request := <-p.RequestVoteRequestChan:
			for i := range p.set {
				go p.set[i].sendRequestVoteRequest(&request)
			}
		}
	}
}
