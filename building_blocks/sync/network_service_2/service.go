// Package main simulates how performant creating a connection with
// Pool to a service that takes time to respond. For every request
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

// create pool of connections to the expensive service
// to help reduce the time required to connect to it
func warmServiceConnCache() *sync.Pool {
	// create a pool of 10 connection to the service
	p := &sync.Pool{
		New: connectToService,
	}
	// instantiate the service and add it to the pool
	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}
	return p
}

func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// get the pool of ready connections
		connPool := warmServiceConnCache()

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
			// get a ready connection to the connectToService and
			// do something with the result of the connection
			srvConn := connPool.Get()
			fmt.Println(conn, "")
			// when you're done with the service connection
			// return it back to the pool to be reused again
			connPool.Put(srvConn)
			conn.Close()
		}
	}()

	return &wg
}
