package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)
func main() {
	l, err := net.Listen("tcp", ":8080") // ListenTCP is needed to set deadline
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	var wg sync.WaitGroup
	sm := make(chan struct{}, 10)
	// closed := make(chan struct{},10)

	ctx,_:= context.WithCancel(context.Background())

	for {
		select {
		case <-ctx.Done():
			log.Println("Context Cancelled.");
		case sm <- struct{}{}:
			wg.Add(1)

			go func() {
				defer func() { <-sm }()
				defer wg.Done()

				conn, err := l.Accept()
				if err != nil {
					log.Println(err)
				}
				handleConnection(conn)
			}()
		}
	}
	wg.Wait()
}


func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte,2048)
	for {
		conn.SetDeadline(time.Now().Add(200 * time.Millisecond))

		n, err := conn.Read(buf)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			} else if err != io.EOF {
				log.Println("read error", err)
				return
			}
		}
		if n == 0 {
			return
		}

		num, err := strconv.Atoi(string(buf[:n]))
		if err != nil {
			fmt.Println(err)
			return
		}
	//		log.Printf("received from %v: %d", conn.RemoteAddr(), num*num)
		conn.Write([]byte(strconv.Itoa(num * num)))
	}
}
