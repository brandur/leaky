package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	clients Server
)

func RunClient() {
	for {
		select {
			case <-clients.AppendEntriesRequestChan:

			case request := <-clients.RequestVoteRequestChan:
				requestData, err := json.Marshal(request)
				if err != nil {
					panic(err)
				}
				broadcast("/request-vote", requestData, handleRequestVoteResponse)
		}
	}
}

func init() {
	clients = Server{
		AppendEntriesRequestChan:  make(chan AppendEntriesRequest, 50),
		AppendEntriesResponseChan: make(chan AppendEntriesResponse, 50),
		RequestVoteRequestChan:    make(chan RequestVoteRequest, 50),
		RequestVoteResponseChan:   make(chan RequestVoteResponse, 50),
	}
}

func broadcast(path string, data []byte, handler func([]byte)) {
	for i := range conf.peerUrls {
		go connectAndSend(conf.peerUrls[i], path, data, handler)
	}
}

func connectAndSend(url string, path string, data []byte, handler func([]byte)) {
    buffer := bytes.NewBuffer(data)
	r, err := http.Post(url + path, "application/json", buffer)
	if err != nil {
		fmt.Printf("send_error peer_url=%v err=\"%v\"", url, err)
		return
	}

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	handler(responseData)
}

func handleRequestVoteResponse(responseData []byte) {
	response := RequestVoteResponse{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		panic(err)
	}
	clients.RequestVoteResponseChan <- response
}
