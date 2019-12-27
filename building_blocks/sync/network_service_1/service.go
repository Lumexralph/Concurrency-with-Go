// Package main simulates creating a connection without Pool to a
// service that takes time to respond. For every request
// made to a server, a new connection is created to the service.
package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()

		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			// make a request to the expensive service
			connectToService()
			fmt.Fprintln(conn, "")
			conn.Close()
		}
	}()

	return &wg
}
