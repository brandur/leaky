package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Url string
}

func newClient(url string) *Client {
	return &Client{
		Url: url,
	}
}

func (c *Client) connectAndSend(path string, data []byte, handler func([]byte)) {
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
