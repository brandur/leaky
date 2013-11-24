package main

import (
	"os"
	"strings"
)

type Conf struct {
	addr string
	peerUrls []string
}

var (
	conf Conf
)

func mustEnv(n string) string {
	v := os.Getenv(n)
	if len(v) == 0 {
		panic("Must set: " + n)
	}
	return v
}

func init() {
	conf.addr = mustEnv("ADDR")
	conf.peerUrls = strings.Split(mustEnv("PEER_URLS"), ",")
}
