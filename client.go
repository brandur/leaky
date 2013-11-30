package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	PeerSet *PeerSet
	Url string
}

func newClient(peerSet *PeerSet, url string) *Client {
	return &Client{
		PeerSet: peerSet,
		Url: url,
	}
}

func (c *Client) buildHandler(handler func(*Client, []byte)) func([]byte) {
	return func(responseData []byte) {
		handler(c, responseData)
	}
}

func (c *Client) send(path string, data []byte, handler func([]byte)) {
	buffer := bytes.NewBuffer(data)
	r, err := http.Post(c.Url+path, "application/json", buffer)
	if err != nil {
		fmt.Printf("send_error peer_url=%v err=\"%v\"", c.Url, err)
		return
	}

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	handler(responseData)
}

func handleAppendEntriesResponse(c *Client, responseData []byte) {
	response := AppendEntriesResponse{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		panic(err)
	}
	c.PeerSet.AppendEntriesResponseChan <- response
}

func handleRequestVoteResponse(c *Client, responseData []byte) {
	response := RequestVoteResponse{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		panic(err)
	}
	c.PeerSet.RequestVoteResponseChan <- response
}

func (c *Client) sendAppendEntriesRequest(request *AppendEntriesRequest) {
	requestData, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	c.send("/append-entries", requestData, c.buildHandler(handleAppendEntriesResponse))
}

func (c *Client) sendRequestVoteRequest(request *RequestVoteRequest) {
	requestData, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	c.send("/request-vote", requestData, c.buildHandler(handleRequestVoteResponse))
}
