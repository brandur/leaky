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
	Term          int  `json:"term"`
	CandidateName Name `json:"candidateName"`
	LastLogIndex  int  `json:"lastLogIndex"`
	LastLogTerm   int  `json:"lastLogTerm"`
}

type RequestVoteResponse struct {
	Term        int  `json:"term"`
	VoteGranted bool `json:"voteGranted"`
}

const (
	ELECTION_TIMEOUT time.Duration = 2 * time.Second

	CANDIDATE State = "candidate"
	FOLLOWER  State = "follower"
	LEADER    State = "leader"
)

var (
	// operational state
	state State
	name  Name

	// persistent state
	currentTerm int = 0
	votedFor    Name

	// volatile state
	commitIndex int = 0
	lastApplied int = 0

	// volatile state for leader
	matchIndex []int
	nextIndex  []int
)

func init() {
	setState(FOLLOWER)
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	name = Name("leaky-" + id.String())
	fmt.Printf("name=%v\n", name)
}

func RunState(server *Server) {
	for {
		select {
		case <-server.AppendEntriesRequestChan:
			fmt.Printf("fuck\n")

		case request := <-server.RequestVoteRequestChan:
			fmt.Printf("vote_requested name=%v\n", request.CandidateName)

			var response RequestVoteResponse

			// grant vote if:
			//     (1) we haven't voted or if we've previously voted for this
			//         client
			//     (2) candidate's log is at least as up-to-date as ours
			if (votedFor == "" || votedFor == request.CandidateName) && request.LastLogTerm >= currentTerm && request.LastLogIndex >= commitIndex {
				response = RequestVoteResponse{Term: currentTerm, VoteGranted: true}
			} else {
				response = RequestVoteResponse{Term: currentTerm, VoteGranted: false}
			}
			server.RequestVoteResponseChan <- response

		// on election timeout, convert to candidate, start election
		case <-time.After(ELECTION_TIMEOUT):
			startElection()
		}
	}
}

func setState(newState State) {
	state = newState
	fmt.Printf("state=%v\n", newState)
}

func setTerm(newTerm int) {
	currentTerm = newTerm
	fmt.Printf("term=%v\n", newTerm)
}

func setVote(vote Name) {
	votedFor = vote
	fmt.Printf("vote=%v\n", vote)
}

func startElection() {
	fmt.Printf("start_election\n")
	setState(CANDIDATE)
	setTerm(currentTerm + 1)
	// vote for self
	setVote(name)
}
