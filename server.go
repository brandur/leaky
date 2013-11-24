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

type ServerError struct {
	Message string `json:"message"`
}

var (
	server Server
)

func RunServer() {
	fmt.Printf("addr=%v\n", conf.addr)
	http.Handle("/", initRouter())
	http.ListenAndServe(conf.addr, nil)
}

func init() {
	server = Server{
		AppendEntriesRequestChan:  make(chan AppendEntriesRequest),
		AppendEntriesResponseChan: make(chan AppendEntriesResponse),
		RequestVoteRequestChan:    make(chan RequestVoteRequest),
		RequestVoteResponseChan:   make(chan RequestVoteResponse),
	}
}

func encodeError(message string) string {
	serverError := ServerError{Message: message}
    responseData, err := json.Marshal(serverError)
	if err != nil {
		panic(err)
	}
	return string(responseData)
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
		http.Error(w, encodeError("Invalid JSON."), http.StatusBadRequest)
		return
	}

	if request.CandidateName == "" {
		http.Error(w, encodeError("Missing field: \"candidateName\"."), http.StatusBadRequest)
		return
	}

	server.RequestVoteRequestChan <- request
	response := <-server.RequestVoteResponseChan

    responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(responseData); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
