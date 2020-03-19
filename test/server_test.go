package test

import (
	"flag"
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {

	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	fmt.Print(addrs)

	var port int
	flag.IntVar(&port, "port", 8001, "Geecache server port")

	fmt.Print(port)

}
