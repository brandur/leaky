package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/pat"
	"io/ioutil"
	"net/http"
)

type Server struct {
	AppendEntriesRequestChan  chan AppendEntriesRequest
	AppendEntriesResponseChan chan AppendEntriesResponse
	RequestVoteRequestChan    chan RequestVoteRequest
	RequestVoteResponseChan   chan RequestVoteResponse
}

var (
	server Server
)

func RunServer() {
	fmt.Printf("port=%v\n", 5100)
	http.Handle("/", initRouter())
	http.ListenAndServe(":5100", nil)
}

func init() {
	server = Server{
		AppendEntriesRequestChan:  make(chan AppendEntriesRequest),
		AppendEntriesResponseChan: make(chan AppendEntriesResponse),
		RequestVoteRequestChan:    make(chan RequestVoteRequest),
		RequestVoteResponseChan:   make(chan RequestVoteResponse),
	}
}

func initRouter() *pat.Router {
	router := pat.New()
	router.Post("/request-vote", requestVoteHandler)
	return router
}

func requestVoteHandler(w http.ResponseWriter, r *http.Request) {
    requestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	request := RequestVoteRequest{}
	if err := json.Unmarshal(requestData, &request); err != nil {
		panic(err)
	}

	server.RequestVoteRequestChan <- request
	response := <-server.RequestVoteResponseChan

    responseData, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	if _, err = w.Write(responseData); err != nil {
		panic(err)
	}
}
