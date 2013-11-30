package main

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"time"
)

type Name string

type State string

type AppendEntriesRequest struct {
}

type AppendEntriesResponse struct {
}

type RequestVoteRequest struct {
	CandidateName Name `json:"candidateName"`
	Term          int  `json:"term"`
	LastLogIndex  int  `json:"lastLogIndex"`
	LastLogTerm   int  `json:"lastLogTerm"`
}

type RequestVoteResponse struct {
	CandidateName Name `json:"candidateName"`
	Term          int  `json:"term"`
	VoteGranted   bool `json:"voteGranted"`
}

const (
	ELECTION_TIMEOUT  time.Duration = 2 * time.Second
	HEARTBEAT_TIMEOUT time.Duration = 1 * time.Second

	CANDIDATE State = "candidate"
	FOLLOWER  State = "follower"
	LEADER    State = "leader"
)

type StateMachine struct {
	peerSet *PeerSet
	server *Server

	// operational state
	state State
	name  Name

	// persistent state
	currentTerm int
	votedFor    Name
	votes       map[Name]bool

	// volatile state
	commitIndex int
	lastApplied int

	// volatile state for leader
	matchIndex []int
	nextIndex  []int
}

func newStateMachine(peerSet *PeerSet, server *Server) *StateMachine {
    s := &StateMachine{}

	s.peerSet = peerSet
	s.server = server
	s.setState(FOLLOWER)

	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	s.name = Name("leaky-" + id.String())
	fmt.Printf("name=%v\n", s.name)
	return s
}

func (s *StateMachine) run() {
	for {
		switch {
		case s.state == CANDIDATE:
			s.runAsCandidate()
		case s.state == FOLLOWER:
			s.runAsFollower()
		case s.state == LEADER:
			s.runAsLeader()
		}
	}
}

func (s *StateMachine) runAsCandidate() {
	for {
		select {
		// peers responding to our vote requests
		case response := <-s.peerSet.RequestVoteResponseChan:
			s.votes[response.CandidateName] = true
			fmt.Printf("num_votes=%v\n", len(s.votes))

			// compare number of votes + 1 for self
			if len(s.votes)+1 > s.peerSet.len()/2 {
				s.setState(LEADER)

				// as we transition into a leadership role, send heartbeat to
				// peers to tell them what's happened
				s.sendHeartbeat()
			}

		// leader requesting an append entries
		case <-s.server.AppendEntriesRequestChan:

		// election timeout; convert to candidate, start election
		case <-time.After(ELECTION_TIMEOUT):
			s.startElection()
		}
	}
}

func (s *StateMachine) runAsFollower() {
	for {
		select {
		// leader requesting an append entries
		case <-s.server.AppendEntriesRequestChan:

		// peers requesting votes
		case request := <-s.server.RequestVoteRequestChan:
			fmt.Printf("vote_requested name=%v\n", request.CandidateName)

			var response RequestVoteResponse

			// grant vote if:
			//     (1) we haven't voted or if we've previously voted for this
			//         client
			//     (2) candidate's log is at least as up-to-date as ours
			if (s.votedFor == "" || s.votedFor == request.CandidateName) && request.LastLogTerm >= s.currentTerm && request.LastLogIndex >= s.commitIndex {
				response = RequestVoteResponse{
					CandidateName: s.name,
					Term:          s.currentTerm,
					VoteGranted:   true,
				}
			} else {
				response = RequestVoteResponse{
					CandidateName: s.name,
					Term:          s.currentTerm,
					VoteGranted:   false,
				}
			}
			s.server.RequestVoteResponseChan <- response

		// election timeout; convert to candidate, start election
		case <-time.After(ELECTION_TIMEOUT):
			s.startElection()
		}
	}
}

func (s *StateMachine) runAsLeader() {
	for {
		select {
		// peers responding to our append entries requests
		case <-s.peerSet.AppendEntriesResponseChan:

		case <-time.After(HEARTBEAT_TIMEOUT):
			s.sendHeartbeat()
		}
	}
}

func (s *StateMachine) sendHeartbeat() {
	request := AppendEntriesRequest{}
	s.peerSet.AppendEntriesRequestChan <- request
}

func (s *StateMachine) setState(state State) {
	s.state = state
	fmt.Printf("state=%v\n", state)

	s.votedFor = ""

	s.votes = make(map[Name]bool)
	fmt.Printf("num_votes=%v\n", len(s.votes))
}

func (s *StateMachine) setTerm(term int) {
	s.currentTerm = term
	fmt.Printf("term=%v\n", term)
}

func (s *StateMachine) setVote(vote Name) {
	s.votedFor = vote
	fmt.Printf("vote=%v\n", vote)
}

func (s *StateMachine) startElection() {
	fmt.Printf("start_election\n")
	s.setState(CANDIDATE)
	s.setTerm(s.currentTerm + 1)

	// vote for self
	s.setVote(s.name)

	request := RequestVoteRequest{
		Term:          s.currentTerm,
		CandidateName: s.name,
		LastLogIndex:  s.commitIndex,
		LastLogTerm:   log[s.commitIndex].term,
	}
	// request vote from other clients
	s.peerSet.RequestVoteRequestChan <- request
}
