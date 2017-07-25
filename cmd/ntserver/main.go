package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
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
	var l, l2, delta int64
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("error at Read.%s\n", err)
			continue
		}
		err = binary.Read(bytes.NewReader(buf[:n]), binary.BigEndian, &l)
		if err != nil {
			fmt.Printf("error at Read.%s\n", err)
			continue
		}
		delta = l - l2
		if delta > 1 {
			ch <- delta
			fmt.Printf("loss: %d\n", delta)
		}
		l2 = l
		i++
		if i%100000 == 0 {
			// fmt.Println(time.Now())
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
