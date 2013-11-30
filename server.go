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

func newServer() *Server {
	return &Server{
		AppendEntriesRequestChan:  make(chan AppendEntriesRequest),
		AppendEntriesResponseChan: make(chan AppendEntriesResponse),
		RequestVoteRequestChan:    make(chan RequestVoteRequest),
		RequestVoteResponseChan:   make(chan RequestVoteResponse),
	}
}

func (s *Server) buildHandler(handler func(*Server, []byte) ([]byte, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		responseData, err := handler(s, requestData)
		if err != nil {
			http.Error(w, s.encodeError(err), http.StatusInternalServerError)
		}

		if _, err = w.Write(responseData); err != nil {
			http.Error(w, s.encodeError(err), http.StatusInternalServerError)
			return
		}
	}
}

func handleAppendEntries(s *Server, requestData []byte) ([]byte, error) {
	request := AppendEntriesRequest{}
	if err := json.Unmarshal(requestData, &request); err != nil {
		return nil, err
	}

	s.AppendEntriesRequestChan <- request
	response := <-s.AppendEntriesResponseChan

	responseData, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}

func handleRequestVote(s *Server, requestData []byte) ([]byte, error) {
	request := RequestVoteRequest{}
	if err := json.Unmarshal(requestData, &request); err != nil {
		return nil, err
	}

	if request.CandidateName == "" {
		return nil, errors.New("Missing field: \"candidateName\".")
	}

	s.RequestVoteRequestChan <- request
	response := <-s.RequestVoteResponseChan

	responseData, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}

func (s *Server) encodeError(err error) string {
	serverError := ServerError{Message: err.Error()}
	responseData, err := json.Marshal(serverError)
	if err != nil {
		panic(err)
	}
	return string(responseData)
}

func (s *Server) initRouter() *pat.Router {
	router := pat.New()
	router.Post("/append-entries", s.buildHandler(handleAppendEntries))
	router.Post("/request-vote", s.buildHandler(handleRequestVote))
	return router
}

func (s *Server) run() {
	fmt.Printf("http=%v\n", conf.http)
	http.Handle("/", s.initRouter())
	http.ListenAndServe(conf.http, nil)
}
