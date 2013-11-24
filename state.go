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
	votes       map[Name]bool

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

func RunState(server *Server, clients *Server) {
	for {
		select {
		// peers responding to our append entries requests
		case <-clients.AppendEntriesResponseChan:

		// peers responding to our vote requests
		case response := <-clients.RequestVoteResponseChan:
			// discard if we're not a candidate
			if state == CANDIDATE {
				votes[response.CandidateName] = true
				fmt.Printf("num_votes=%v\n", len(votes))

				// compare number of votes + 1 for self
				if len(votes) + 1 > len(conf.peerUrls) / 2 {
					setState(LEADER)
				}
			}

		// leader requesting an append entries
		case <-server.AppendEntriesRequestChan:

		// peers requesting votes
		case request := <-server.RequestVoteRequestChan:
			fmt.Printf("vote_requested name=%v\n", request.CandidateName)

			var response RequestVoteResponse

			// grant vote if:
			//     (1) we haven't voted or if we've previously voted for this
			//         client
			//     (2) candidate's log is at least as up-to-date as ours
			if (votedFor == "" || votedFor == request.CandidateName) && request.LastLogTerm >= currentTerm && request.LastLogIndex >= commitIndex {
				response = RequestVoteResponse{
					CandidateName: name,
					Term: currentTerm,
					VoteGranted: true,
				}
			} else {
				response = RequestVoteResponse{
					CandidateName: name,
					Term: currentTerm,
					VoteGranted: false,
				}
			}
			server.RequestVoteResponseChan <- response

		// election timeout; convert to candidate, start election
		case <-time.After(ELECTION_TIMEOUT):
			startElection()
		}
	}
}

func setState(newState State) {
	state = newState
	fmt.Printf("state=%v\n", newState)

	votedFor = ""

	votes = make(map[Name]bool)
	fmt.Printf("num_votes=%v\n", len(votes))
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

	request := RequestVoteRequest{
		Term: currentTerm,
		CandidateName: name,
		LastLogIndex: commitIndex,
		LastLogTerm: log[commitIndex].term,
	}
	// request vote from other clients
	clients.RequestVoteRequestChan <- request
}
