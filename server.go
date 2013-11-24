package main

import (
	"encoding/json"
	"errors"
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

func buildHandler(handler func ([]byte) ([]byte, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		responseData, err := handler(requestData)
		if err != nil {
			http.Error(w, encodeError(err), http.StatusInternalServerError)
		}

		if _, err = w.Write(responseData); err != nil {
			http.Error(w, encodeError(err), http.StatusInternalServerError)
			return
		}
	}
}

func handleAppendEntries(requestData []byte) ([]byte, error) {
	request := AppendEntriesRequest{}
	if err := json.Unmarshal(requestData, &request); err != nil {
		return nil, err
	}

	server.AppendEntriesRequestChan <- request
	response := <-server.AppendEntriesResponseChan

    responseData, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}

func handleRequestVote(requestData []byte) ([]byte, error) {
	request := RequestVoteRequest{}
	if err := json.Unmarshal(requestData, &request); err != nil {
		return nil, err
	}

	if request.CandidateName == "" {
		return nil, errors.New("Missing field: \"candidateName\".")
	}

	server.RequestVoteRequestChan <- request
	response := <-server.RequestVoteResponseChan

    responseData, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}

func encodeError(err error) string {
	serverError := ServerError{Message: err.Error()}
    responseData, err := json.Marshal(serverError)
	if err != nil {
		panic(err)
	}
	return string(responseData)
}

func initRouter() *pat.Router {
	router := pat.New()
	router.Post("/append-entries", buildHandler(handleAppendEntries))
	router.Post("/request-vote", buildHandler(handleRequestVote))
	return router
}
