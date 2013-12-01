package main

import (
	"os"
)

type Conf struct {
	http    string
	peerUrl string
}

var (
	conf Conf
)

func envRequired(n string) string {
	v := os.Getenv(n)
	if len(v) == 0 {
		panic("Must set: " + n)
	}
	return v
}

func init() {
	conf.http = os.Getenv("HTTP")
	if conf.http == "" {
		conf.http = ":3131"
	}
	conf.peerUrl = os.Getenv("PEER_URL")
}
