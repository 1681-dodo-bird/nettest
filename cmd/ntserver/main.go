package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/1681-dodo-bird/nettest/msg"
)

func server(ch chan int64, ch2 chan int) {
	serverAddr, err := net.ResolveUDPAddr("udp", ":12299")
	if err != nil {
		fmt.Printf("Cannot listen:%s\n", err)
		return
	}

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Printf("Cannot listen:%s\n", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	i := 0
	var l2, delta int64
	m := msg.Message{}
	for {
		n, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("error at Read.%s\n", err)
			continue
		}
		err = binary.Read(bytes.NewReader(buf[:n]), binary.BigEndian, &m)
		if err != nil {
			fmt.Printf("error at Read.%s\n", err)
			continue
		}
		delta = m.I - l2
		if delta > 1 {
			ch <- delta
		}
		l2 = m.I

		conn.WriteTo(buf[:n], remote)
		i++
		if i%100000 == 0 {
			i = 0
			ch2 <- 0
		}
	}

}

func main() {
	ch := make(chan int64)
	ch2 := make(chan int)
	go server(ch, ch2)

	for {
		select {
		case v := <-ch:
			fmt.Printf("loss %v\n", v)
		case <-ch2:
			fmt.Println(time.Now())
		}
	}
}
