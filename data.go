package main

type Name string

type AppendEntriesRequest struct {
}

type AppendEntriesResponse struct {
}

type JoinRequest struct {
	Name Name `json:"name"`
	Url string `json:"url"`
}

type JoinResponse struct {
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
