package main

import (
	"os"
	"strings"
)

type Conf struct {
	peerUrls []string
}

var (
	conf Conf
)

func mustEnv(n string) string {
	v := os.Getenv(n)
//	if len(v) == 0 {
//		panic("Must set: " + n)
//	}
	return v
}

func init() {
	conf.peerUrls = strings.Split(mustEnv("PEER_URLS"), ",")
}
